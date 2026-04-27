package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/portfolio/backend/internal/service"
)

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
}

func buildHandler(t *testing.T, upstream *httptest.Server) *AnalyticsHandler {
	t.Helper()

	svc, err := service.NewAnalyticsService(
		upstream.URL,
		"test-key",
		"00000000-0000-0000-0000-000000000001",
		"dev",
		time.Minute,
		quietLogger(),
	)
	if err != nil {
		t.Fatalf("service construct: %v", err)
	}
	return NewAnalyticsHandler(svc, quietLogger())
}

func buildDisabledHandler(t *testing.T) *AnalyticsHandler {
	t.Helper()
	svc, err := service.NewAnalyticsService("", "", "", "", 0, quietLogger())
	if err != nil {
		t.Fatalf("service construct: %v", err)
	}
	if !svc.Disabled() {
		t.Fatalf("expected disabled service")
	}
	return NewAnalyticsHandler(svc, quietLogger())
}

type testEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func decodeEnvelope(t *testing.T, rec *httptest.ResponseRecorder) testEnvelope {
	t.Helper()
	var env testEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v — body=%q", err, rec.Body.String())
	}
	return env
}

func newStatsUpstream() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{
          "pageviews": 500,
          "visitors":  100,
          "visits":    100,
          "bounces":   20,
          "totaltime": 6000,
          "comparison": {
            "pageviews": 400,
            "visitors":  0,
            "visits":    0,
            "bounces":   0,
            "totaltime": 0
          }
        }`)
	}))
}

func TestHandler_Summary_Happy(t *testing.T) {
	t.Parallel()

	up := newStatsUpstream()
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary?range=30d", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type: %q", ct)
	}
	env := decodeEnvelope(t, rec)
	if !env.Success {
		t.Fatalf("success=false; err=%q", env.Error)
	}
	if env.Error != "" {
		t.Errorf("expected empty error, got %q", env.Error)
	}
	var data map[string]any
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	if data["range"] != "30d" {
		t.Errorf("range in payload: %v", data["range"])
	}
}

func TestHandler_Summary_InvalidRange(t *testing.T) {
	t.Parallel()

	up := newStatsUpstream()
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary?range=90d", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if env.Success {
		t.Errorf("success=true on invalid range")
	}
	if env.Error != "invalid range" {
		t.Errorf("error msg: got %q, want 'invalid range'", env.Error)
	}
}

func TestHandler_Summary_Disabled503(t *testing.T) {
	t.Parallel()

	h := buildDisabledHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if env.Success {
		t.Errorf("success=true on disabled")
	}
	if env.Error != "analytics_not_configured" {
		t.Errorf("error: got %q, want analytics_not_configured", env.Error)
	}
}

func TestHandler_Summary_Upstream5xx(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary?range=7d", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503 (5xx→503); body=%s", rec.Code, rec.Body.String())
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "analytics unavailable" {
		t.Errorf("error: got %q, want 'analytics unavailable'", env.Error)
	}
}

func TestHandler_Summary_Upstream4xxNon401(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "no such site", http.StatusNotFound)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary?range=7d", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status: got %d, want 502 (4xx→502); body=%s", rec.Code, rec.Body.String())
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "analytics upstream error" {
		t.Errorf("error: got %q, want 'analytics upstream error'", env.Error)
	}
}

func TestHandler_Summary_Upstream401MaskedTo503(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary?range=7d", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503 (mask 401); body=%s", rec.Code, rec.Body.String())
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "analytics unavailable" {
		t.Errorf("error: got %q, want 'analytics unavailable'", env.Error)
	}
	if strings.Contains(rec.Body.String(), "nope") ||
		strings.Contains(rec.Body.String(), "401") {
		t.Errorf("upstream detail leaked in body: %s", rec.Body.String())
	}
}

func TestHandler_TopPages_Happy(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[{"x":"/ru","y":10},{"x":"/ru/about","y":5}]`)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/top-pages?range=7d&limit=5", nil)
	rec := httptest.NewRecorder()
	h.TopPages(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if !env.Success {
		t.Fatalf("success=false; err=%q", env.Error)
	}
	var arr []map[string]any
	if err := json.Unmarshal(env.Data, &arr); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("len: got %d, want 2", len(arr))
	}
}

func TestHandler_TopPages_InvalidLimitNonInt(t *testing.T) {
	t.Parallel()

	up := newStatsUpstream()
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/top-pages?limit=abc", nil)
	rec := httptest.NewRecorder()
	h.TopPages(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "invalid limit" {
		t.Errorf("error: got %q", env.Error)
	}
}

func TestHandler_TopPages_InvalidLimitOutOfRange(t *testing.T) {
	t.Parallel()

	up := newStatsUpstream()
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	for _, raw := range []string{"0", "-5", "101", "99999"} {
		t.Run("limit="+raw, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/api/analytics/top-pages?limit="+raw, nil)
			rec := httptest.NewRecorder()
			h.TopPages(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("limit=%s: got %d, want 400", raw, rec.Code)
			}
		})
	}
}

func TestHandler_TopReferrers_Happy(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[{"x":"google.com","y":10},{"x":"","y":5}]`)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/top-referrers?range=30d", nil)
	rec := httptest.NewRecorder()
	h.TopReferrers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	var env testEnvelope
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if !env.Success {
		t.Fatalf("success=false: %q", env.Error)
	}
}

func TestHandler_TopCountries_Disabled(t *testing.T) {
	t.Parallel()

	h := buildDisabledHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/top-countries", nil)
	rec := httptest.NewRecorder()
	h.TopCountries(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503", rec.Code)
	}
}

func TestHandler_Timeseries_Happy(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{
          "pageviews":[{"x":"2026-03-25","y":100}],
          "sessions": [{"x":"2026-03-25","y":40}]
        }`)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/timeseries?range=7d&metric=pageviews", nil)
	rec := httptest.NewRecorder()
	h.Timeseries(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	env := decodeEnvelope(t, rec)
	if !env.Success {
		t.Fatalf("success=false: %q", env.Error)
	}
}

func TestHandler_Timeseries_InvalidMetric(t *testing.T) {
	t.Parallel()

	up := newStatsUpstream()
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/timeseries?metric=clicks", nil)
	rec := httptest.NewRecorder()
	h.Timeseries(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "invalid metric" {
		t.Errorf("error: got %q", env.Error)
	}
}

func newMetricsUpstream(t *testing.T, recordType *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if recordType != nil {
			*recordType = r.URL.Query().Get("type")
		}
		_, _ = io.WriteString(w, `[{"x":"instagram","y":10},{"x":"","y":5}]`)
	}))
}

func TestHandler_TopUTM_MissingType(t *testing.T) {
	t.Parallel()

	up := newMetricsUpstream(t, nil)
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/top-utm?range=7d&limit=10", nil)
	rec := httptest.NewRecorder()
	h.TopUTM(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if env.Success {
		t.Errorf("success=true on missing type")
	}
	if env.Error != "utm type required" {
		t.Errorf("error: got %q, want 'utm type required'", env.Error)
	}
}

func TestHandler_TopUTM_RestrictedTypes(t *testing.T) {
	t.Parallel()

	up := newMetricsUpstream(t, nil)
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	for _, typeStr := range []string{"term", "content", "foo", "utm_source", "Source", " ", "sourceX"} {
		typeStr := typeStr
		t.Run("reject_"+typeStr, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet,
				"/api/analytics/top-utm?type="+url.QueryEscape(typeStr)+"&range=7d", nil)
			rec := httptest.NewRecorder()
			h.TopUTM(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("type=%q: got %d, want 400; body=%s", typeStr, rec.Code, rec.Body.String())
			}
			env := decodeEnvelope(t, rec)
			if env.Error != "invalid utm type" {
				t.Errorf("type=%q: error got %q, want 'invalid utm type'", typeStr, env.Error)
			}
		})
	}
}

func TestHandler_TopUTM_ValidTypes(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"source":   "utmSource",
		"medium":   "utmMedium",
		"campaign": "utmCampaign",
	}

	for short, want := range cases {
		short, want := short, want
		t.Run("accept_"+short, func(t *testing.T) {
			t.Parallel()

			var gotType string
			up := newMetricsUpstream(t, &gotType)
			t.Cleanup(up.Close)
			h := buildHandler(t, up)

			req := httptest.NewRequest(http.MethodGet,
				"/api/analytics/top-utm?type="+short+"&range=7d&limit=5", nil)
			rec := httptest.NewRecorder()
			h.TopUTM(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("type=%s: got %d, want 200; body=%s", short, rec.Code, rec.Body.String())
			}
			if gotType != want {
				t.Errorf("upstream type: got %q, want %q", gotType, want)
			}

			env := decodeEnvelope(t, rec)
			if !env.Success {
				t.Fatalf("success=false: %q", env.Error)
			}
			var arr []map[string]any
			if err := json.Unmarshal(env.Data, &arr); err != nil {
				t.Fatalf("decode data: %v", err)
			}
			if len(arr) != 2 {
				t.Fatalf("len: got %d, want 2", len(arr))
			}
			if arr[0]["value"] != "instagram" {
				t.Errorf("[0].value: got %v, want instagram", arr[0]["value"])
			}
			if arr[1]["value"] != "(none)" {
				t.Errorf("[1].value (empty→none): got %v", arr[1]["value"])
			}
		})
	}
}

func TestHandler_TopUTM_InvalidLimit(t *testing.T) {
	t.Parallel()

	up := newMetricsUpstream(t, nil)
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	for _, raw := range []string{"abc", "0", "-5", "101"} {
		raw := raw
		t.Run("limit="+raw, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet,
				"/api/analytics/top-utm?type=source&limit="+raw, nil)
			rec := httptest.NewRecorder()
			h.TopUTM(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("limit=%s: got %d, want 400", raw, rec.Code)
			}
		})
	}
}

func TestHandler_TopUTM_Disabled(t *testing.T) {
	t.Parallel()

	h := buildDisabledHandler(t)
	req := httptest.NewRequest(http.MethodGet,
		"/api/analytics/top-utm?type=source&range=7d", nil)
	rec := httptest.NewRecorder()
	h.TopUTM(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503", rec.Code)
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "analytics_not_configured" {
		t.Errorf("error: got %q", env.Error)
	}
}

func TestHandler_TopUTM_Upstream4xx(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad req", http.StatusBadRequest)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet,
		"/api/analytics/top-utm?type=source&range=7d", nil)
	rec := httptest.NewRecorder()
	h.TopUTM(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status: got %d, want 502; body=%s", rec.Code, rec.Body.String())
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "analytics upstream error" {
		t.Errorf("error: got %q", env.Error)
	}
}

func TestHandler_TopUTM_Upstream5xx(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet,
		"/api/analytics/top-utm?type=source&range=7d", nil)
	rec := httptest.NewRecorder()
	h.TopUTM(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503; body=%s", rec.Code, rec.Body.String())
	}
	env := decodeEnvelope(t, rec)
	if env.Error != "analytics unavailable" {
		t.Errorf("error: got %q", env.Error)
	}
}

func TestHandler_EnvelopeShape(t *testing.T) {
	t.Parallel()

	up := newStatsUpstream()
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := raw["success"]; !ok {
		t.Errorf("missing 'success' key")
	}
	if _, ok := raw["data"]; !ok {
		t.Errorf("missing 'data' key")
	}
	if _, ok := raw["error"]; ok {
		t.Errorf("'error' key present on success")
	}
}

func TestHandler_NoUpstreamURLInBody(t *testing.T) {
	t.Parallel()

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "internal detail", http.StatusInternalServerError)
	}))
	t.Cleanup(up.Close)
	h := buildHandler(t, up)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.Summary(rec, req)

	body := rec.Body.String()
	forbidden := []string{
		"test-key",
		"Bearer",
		"00000000-0000-0000-0000-000000000001",
		"internal detail",
	}
	if u, err := url.Parse(up.URL); err == nil && u.Host != "" {
		forbidden = append(forbidden, u.Host)
	}
	for _, bad := range forbidden {
		if bad != "" && strings.Contains(body, bad) {
			t.Errorf("body leaks %q: %s", bad, body)
		}
	}
}
