package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newTestService(t *testing.T, ts *httptest.Server, ttl time.Duration) *AnalyticsService {
	t.Helper()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("parse httptest url: %v", err)
	}

	s := &AnalyticsService{
		client: &http.Client{
			Timeout: 2 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		apiURL:    u,
		apiKey:    "test-key",
		websiteID: "00000000-0000-0000-0000-000000000001",
		cache:     newTTLCache(),
		ttl:       ttl,
		logger:    silentLogger(),
	}
	return s
}

func umamiStatsBody(pv, prev, visitors, visits, bounces, totaltime int64) string {
	return fmt.Sprintf(`{
  "pageviews": {"value": %d, "prev": %d},
  "visitors":  {"value": %d, "prev": 0},
  "visits":    {"value": %d, "prev": 0},
  "bounces":   {"value": %d, "prev": 0},
  "totaltime": {"value": %d, "prev": 0}
}`, pv, prev, visitors, visits, bounces, totaltime)
}

func TestNewAnalyticsService_DisabledWhenEnvsMissing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		apiURL    string
		apiKey    string
		websiteID string
	}{
		{"all empty", "", "", ""},
		{"only url", "http://umami", "", ""},
		{"only url+key", "http://umami", "k", ""},
		{"missing url", "", "k", "w"},
		{"missing key", "http://umami", "", "w"},
		{"missing website", "http://umami", "k", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, err := NewAnalyticsService(tc.apiURL, tc.apiKey, tc.websiteID, "dev", 0, silentLogger())
			if err != nil {
				t.Fatalf("unexpected constructor error: %v", err)
			}
			if s == nil {
				t.Fatalf("expected non-nil service in disabled mode")
			}
			if !s.Disabled() {
				t.Fatalf("expected Disabled()=true")
			}

			ctx := context.Background()
			if _, err := s.Summary(ctx, "7d"); !errors.Is(err, ErrServiceDisabled) {
				t.Errorf("Summary: got %v, want ErrServiceDisabled", err)
			}
			if _, err := s.TopPages(ctx, "7d", 10); !errors.Is(err, ErrServiceDisabled) {
				t.Errorf("TopPages: got %v, want ErrServiceDisabled", err)
			}
			if _, err := s.TopReferrers(ctx, "7d", 10); !errors.Is(err, ErrServiceDisabled) {
				t.Errorf("TopReferrers: got %v, want ErrServiceDisabled", err)
			}
			if _, err := s.TopCountries(ctx, "7d", 10); !errors.Is(err, ErrServiceDisabled) {
				t.Errorf("TopCountries: got %v, want ErrServiceDisabled", err)
			}
			if _, err := s.Timeseries(ctx, "7d", "pageviews"); !errors.Is(err, ErrServiceDisabled) {
				t.Errorf("Timeseries: got %v, want ErrServiceDisabled", err)
			}
		})
	}
}

func TestNewAnalyticsService_InvalidUUID(t *testing.T) {
	t.Parallel()

	bad := []string{
		"not-a-uuid",
		"00000000-0000-0000-0000",
		"00000000-0000-0000-0000-00000000000G",
		"00000000-0000-0000-0000-0000000000001",
		"00000000-0000-0000-0000-000000000001 ",
		"../../etc/passwd",
		"../../../websites/other/stats?",
		"00000000-0000-0000-0000-000000000001/../other",
	}

	for _, id := range bad {
		id := id
		t.Run(fmt.Sprintf("id=%q", id), func(t *testing.T) {
			t.Parallel()
			s, err := NewAnalyticsService(
				"http://umami:3000", "k", id, "prod", 0, silentLogger(),
			)
			if err == nil {
				t.Fatalf("expected error for websiteID=%q, got nil (service=%v)", id, s)
			}
			if s != nil {
				t.Errorf("expected nil service on invalid UUID, got %v", s)
			}
		})
	}
}

func TestNewAnalyticsService_ValidUUID(t *testing.T) {
	t.Parallel()

	good := []string{
		"12345678-1234-1234-1234-123456789012",
		"00000000-0000-0000-0000-000000000001",
		"abcdef01-2345-6789-abcd-ef0123456789",
	}

	for _, id := range good {
		id := id
		t.Run(fmt.Sprintf("id=%s", id), func(t *testing.T) {
			t.Parallel()
			s, err := NewAnalyticsService(
				"http://umami:3000", "k", id, "prod", 0, silentLogger(),
			)
			if err != nil {
				t.Fatalf("unexpected error for valid UUID %q: %v", id, err)
			}
			if s == nil {
				t.Fatal("expected non-nil service for valid UUID")
			}
			if s.Disabled() {
				t.Error("expected Disabled()=false for fully configured service")
			}
		})
	}
}

func TestNewAnalyticsService_URLValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		apiURL  string
		env     string
		wantErr bool
	}{
		{"valid internal docker dns", "http://umami:3000", "prod", false},
		{"valid external fqdn in prod", "https://umami.lavier.tech", "prod", false},
		{"valid dev localhost", "http://localhost:3000", "dev", false},
		{"localhost rejected in prod", "http://localhost:3000", "prod", true},
		{"public host accepted as operator-trusted", "https://analytics.example.com", "dev", false},
		{"bad scheme", "file:///etc/passwd", "dev", true},
		{"aws metadata rejected", "http://169.254.169.254", "dev", true},
		{"gcp metadata hostname rejected", "http://metadata.google.internal", "prod", true},
		{"link-local ipv4 rejected", "http://169.254.1.2", "dev", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewAnalyticsService(
				tc.apiURL, "k",
				"00000000-0000-0000-0000-000000000001",
				tc.env, 0, silentLogger(),
			)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestService_InputValidation(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Errorf("upstream hit during validation-only test")
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	ctx := context.Background()

	t.Run("invalid range", func(t *testing.T) {
		t.Parallel()
		if _, err := s.Summary(ctx, "90d"); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("Summary: got %v, want ErrInvalidRange", err)
		}
		if _, err := s.TopPages(ctx, "1y", 10); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("TopPages: got %v, want ErrInvalidRange", err)
		}
		if _, err := s.TopReferrers(ctx, "", 10); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("TopReferrers empty: got %v, want ErrInvalidRange", err)
		}
		if _, err := s.TopCountries(ctx, "bogus", 10); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("TopCountries: got %v, want ErrInvalidRange", err)
		}
		if _, err := s.Timeseries(ctx, "weird", "pageviews"); !errors.Is(err, ErrInvalidRange) {
			t.Errorf("Timeseries: got %v, want ErrInvalidRange", err)
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		t.Parallel()
		limits := []int{0, -1, 101, 99999}
		for _, l := range limits {
			if _, err := s.TopPages(ctx, "7d", l); !errors.Is(err, ErrInvalidLimit) {
				t.Errorf("limit=%d: got %v, want ErrInvalidLimit", l, err)
			}
		}
	})

	t.Run("invalid metric", func(t *testing.T) {
		t.Parallel()
		bads := []string{"", "clicks", "Pageviews", "unique"}
		for _, m := range bads {
			if _, err := s.Timeseries(ctx, "7d", m); !errors.Is(err, ErrInvalidMetric) {
				t.Errorf("metric=%q: got %v, want ErrInvalidMetric", m, err)
			}
		}
	})
}

func TestSummary_HappyPath(t *testing.T) {
	t.Parallel()

	var gotAuth string
	var gotPath string
	var hits int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, umamiStatsBody(1000, 800, 456, 500, 200, 60000))
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	sum, err := s.Summary(context.Background(), "30d")
	if err != nil {
		t.Fatalf("Summary error: %v", err)
	}

	if gotAuth != "Bearer test-key" {
		t.Errorf("Authorization header mismatch: %q", gotAuth)
	}
	if !strings.HasSuffix(gotPath, "/api/websites/00000000-0000-0000-0000-000000000001/stats") {
		t.Errorf("upstream path mismatch: %q", gotPath)
	}
	if sum.Range != "30d" {
		t.Errorf("Range: got %q, want 30d", sum.Range)
	}
	if sum.Pageviews != 1000 {
		t.Errorf("Pageviews: got %d, want 1000", sum.Pageviews)
	}
	if sum.PreviousPageviews != 800 {
		t.Errorf("PreviousPageviews: got %d, want 800", sum.PreviousPageviews)
	}
	if sum.UniqueVisitors != 456 {
		t.Errorf("UniqueVisitors: got %d, want 456", sum.UniqueVisitors)
	}
	if sum.BounceRate < 0.39 || sum.BounceRate > 0.41 {
		t.Errorf("BounceRate: got %v, want ~0.4", sum.BounceRate)
	}
	if sum.AvgSessionSeconds != 120 {
		t.Errorf("AvgSessionSeconds: got %d, want 120", sum.AvgSessionSeconds)
	}
	if sum.DeltaPageviews < 0.24 || sum.DeltaPageviews > 0.26 {
		t.Errorf("DeltaPageviews: got %v, want ~0.25", sum.DeltaPageviews)
	}

	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("upstream hits: got %d, want 1", got)
	}
}

func TestSummary_BounceRateAndDeltaClamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantBounce float64
		wantDelta  float64
	}{
		{
			name:       "zero visits yields zero bounce and zero avg",
			body:       umamiStatsBody(10, 10, 5, 0, 0, 0),
			wantBounce: 0,
			wantDelta:  0,
		},
		{
			name:       "bounces > visits clamped to 1",
			body:       umamiStatsBody(100, 100, 50, 10, 50, 0),
			wantBounce: 1,
			wantDelta:  0,
		},
		{
			name:       "prev zero + current positive yields delta=10",
			body:       umamiStatsBody(500, 0, 100, 100, 10, 1000),
			wantBounce: 0.1,
			wantDelta:  10,
		},
		{
			name:       "prev zero + current zero yields delta=0",
			body:       umamiStatsBody(0, 0, 0, 0, 0, 0),
			wantBounce: 0,
			wantDelta:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body := tc.body
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, body)
			}))
			t.Cleanup(ts.Close)

			s := newTestService(t, ts, time.Minute)
			sum, err := s.Summary(context.Background(), "7d")
			if err != nil {
				t.Fatalf("Summary: %v", err)
			}
			if sum.BounceRate+0.01 < tc.wantBounce || sum.BounceRate-0.01 > tc.wantBounce {
				t.Errorf("BounceRate: got %v, want ~%v", sum.BounceRate, tc.wantBounce)
			}
			if sum.DeltaPageviews+0.01 < tc.wantDelta || sum.DeltaPageviews-0.01 > tc.wantDelta {
				t.Errorf("DeltaPageviews: got %v, want ~%v", sum.DeltaPageviews, tc.wantDelta)
			}
		})
	}
}

func TestSummary_Upstream5xx(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	_, err := s.Summary(context.Background(), "7d")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrUpstream) {
		t.Errorf("errors.Is(err, ErrUpstream) = false; err=%v", err)
	}
	var up *UpstreamError
	if !errors.As(err, &up) {
		t.Fatalf("errors.As to *UpstreamError failed; err=%v", err)
	}
	if up.Status != http.StatusInternalServerError {
		t.Errorf("UpstreamError.Status: got %d, want 500", up.Status)
	}
}

func TestSummary_UpstreamMalformedJSON(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"pageviews": `)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on malformed JSON: %v", r)
		}
	}()

	_, err := s.Summary(context.Background(), "7d")
	if err == nil {
		t.Fatal("expected error on malformed JSON, got nil")
	}
}

func TestSummary_ContextCancelled(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		_, err := s.Summary(ctx, "7d")
		errCh <- err
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error after cancel, got nil")
		}
		if !errors.Is(err, context.Canceled) && !errors.Is(err, ErrUpstream) {
			t.Errorf("expected context.Canceled or ErrUpstream, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Summary did not return after ctx cancel")
	}
}

func TestTopPages_Mapping(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("type"); got != "url" {
			t.Errorf("type query: got %q, want url", got)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit query: got %q, want 5", got)
		}
		_, _ = io.WriteString(w, `[
          {"x": "https://lavier.tech/ru/projects?a=1#b", "y": 321},
          {"x": "/ru/about", "y": 210},
          {"x": "", "y": 1}
        ]`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	got, err := s.TopPages(context.Background(), "7d", 5)
	if err != nil {
		t.Fatalf("TopPages: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len: got %d, want 3", len(got))
	}
	if got[0].Path != "/ru/projects" {
		t.Errorf("[0].Path: got %q, want /ru/projects", got[0].Path)
	}
	if got[0].Views != 321 {
		t.Errorf("[0].Views: got %d, want 321", got[0].Views)
	}
	if got[0].Uniques != 0 {
		t.Errorf("[0].Uniques: got %d, want 0 (not provided)", got[0].Uniques)
	}
	if got[1].Path != "/ru/about" {
		t.Errorf("[1].Path: got %q, want /ru/about", got[1].Path)
	}
	if got[2].Path != "/" {
		t.Errorf("[2].Path empty→/: got %q", got[2].Path)
	}
}

func TestTopReferrers_EmptyNormalized(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[
          {"x": "google.com", "y": 10},
          {"x": "",           "y": 40},
          {"x": "  ",         "y": 2}
        ]`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	got, err := s.TopReferrers(context.Background(), "30d", 10)
	if err != nil {
		t.Fatalf("TopReferrers: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len: got %d, want 3", len(got))
	}
	if got[0].Referrer != "google.com" {
		t.Errorf("[0]: got %q", got[0].Referrer)
	}
	if got[1].Referrer != "(direct)" {
		t.Errorf("empty → (direct), got %q", got[1].Referrer)
	}
	if got[2].Referrer != "(direct)" {
		t.Errorf("whitespace → (direct), got %q", got[2].Referrer)
	}
}

func TestTopCountries_Normalized(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[
          {"x": "ru", "y": 8000},
          {"x": "US", "y": 1000},
          {"x": "",   "y": 5},
          {"x": "XXX","y": 3},
          {"x": "r1", "y": 2}
        ]`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	got, err := s.TopCountries(context.Background(), "30d", 10)
	if err != nil {
		t.Fatalf("TopCountries: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("len: got %d, want 5", len(got))
	}
	if got[0].Country != "RU" {
		t.Errorf("[0] lowercase→upper: got %q", got[0].Country)
	}
	if got[1].Country != "US" {
		t.Errorf("[1]: got %q", got[1].Country)
	}
	if got[2].Country != "ZZ" {
		t.Errorf("empty→ZZ: got %q", got[2].Country)
	}
	if got[3].Country != "ZZ" {
		t.Errorf("three-letter→ZZ: got %q", got[3].Country)
	}
	if got[4].Country != "R1" {
		t.Errorf("two-char mixed: got %q", got[4].Country)
	}
}

func TestTimeseries_MetricAndDateNormalization(t *testing.T) {
	t.Parallel()

	body := `{
        "pageviews": [
          {"x": "2026-03-25 00:00:00", "y": 412},
          {"x": "2026-03-26",           "y": 389},
          {"x": "2026-03-24T00:00:00Z", "y": 100}
        ],
        "sessions": [
          {"x": "2026-03-25", "y": 40},
          {"x": "2026-03-26", "y": 38}
        ]
      }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("unit"); got != "day" {
			t.Errorf("unit: got %q, want day", got)
		}
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	t.Run("pageviews metric + sort", func(t *testing.T) {
		got, err := s.Timeseries(context.Background(), "30d", "pageviews")
		if err != nil {
			t.Fatalf("Timeseries: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("len: got %d, want 3", len(got))
		}
		if got[0].Date != "2026-03-24" || got[1].Date != "2026-03-25" || got[2].Date != "2026-03-26" {
			t.Errorf("sort: got %+v", got)
		}
		if got[1].Value != 412 {
			t.Errorf("value mapping: got %d, want 412", got[1].Value)
		}
	})

	t.Run("visitors metric uses sessions", func(t *testing.T) {
		got, err := s.Timeseries(context.Background(), "30d", "visitors")
		if err != nil {
			t.Fatalf("Timeseries: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("len: got %d, want 2", len(got))
		}
		if got[0].Value != 40 || got[1].Value != 38 {
			t.Errorf("values: got %+v", got)
		}
	})
}

func TestCache_HitsWithinTTL(t *testing.T) {
	t.Parallel()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = io.WriteString(w, umamiStatsBody(10, 5, 3, 3, 1, 60))
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	for i := 0; i < 5; i++ {
		if _, err := s.Summary(context.Background(), "7d"); err != nil {
			t.Fatalf("Summary #%d: %v", i, err)
		}
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("upstream hits: got %d, want 1 (cached)", got)
	}
}

func TestCache_ExpiresAfterTTL(t *testing.T) {
	t.Parallel()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = io.WriteString(w, umamiStatsBody(1, 1, 1, 1, 0, 1))
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, 20*time.Millisecond)

	if _, err := s.Summary(context.Background(), "7d"); err != nil {
		t.Fatalf("call 1: %v", err)
	}
	time.Sleep(40 * time.Millisecond)
	if _, err := s.Summary(context.Background(), "7d"); err != nil {
		t.Fatalf("call 2: %v", err)
	}

	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Errorf("upstream hits: got %d, want 2", got)
	}
}

func TestCache_ErrorsNotCached(t *testing.T) {
	t.Parallel()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		if n == 1 {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		_, _ = io.WriteString(w, umamiStatsBody(7, 7, 2, 2, 0, 4))
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	if _, err := s.Summary(context.Background(), "7d"); err == nil {
		t.Fatal("call 1: expected error, got nil")
	}
	if _, err := s.Summary(context.Background(), "7d"); err != nil {
		t.Fatalf("call 2: %v", err)
	}
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Errorf("upstream hits: got %d, want 2", got)
	}
}

func TestCache_InFlightDedup(t *testing.T) {
	t.Parallel()

	var hits int32
	release := make(chan struct{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		select {
		case <-release:
		case <-r.Context().Done():
			return
		}
		_, _ = io.WriteString(w, umamiStatsBody(1, 1, 1, 1, 0, 1))
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	errs := make([]error, N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := s.Summary(context.Background(), "7d")
			errs[i] = err
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	close(release)
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d error: %v", i, err)
		}
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("upstream hits under dedup: got %d, want 1", got)
	}
}

func TestTopPages_UpstreamGarbageNoPanic(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `not a json`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on garbage: %v", r)
		}
	}()

	if _, err := s.TopPages(context.Background(), "7d", 10); err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestTopPages_JSONNumberDecimal(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[{"x":"/","y":10},{"x":"/a","y":"5"}]`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	got, err := s.TopPages(context.Background(), "7d", 10)
	if err != nil {
		return
	}
	if len(got) != 2 {
		t.Fatalf("len: got %d, want 2", len(got))
	}
	if got[0].Views != 10 {
		t.Errorf("[0].Views: got %d, want 10", got[0].Views)
	}
}

func TestValidateUtmType(t *testing.T) {
	t.Parallel()

	valid := map[string]UtmType{
		"source":   UtmSource,
		"medium":   UtmMedium,
		"campaign": UtmCampaign,
		"term":     UtmTerm,
		"content":  UtmContent,
	}
	for in, want := range valid {
		in, want := in, want
		t.Run("valid_"+in, func(t *testing.T) {
			t.Parallel()
			got, err := ValidateUtmType(in)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", in, err)
			}
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}

	invalid := []string{"", "Source", "utm_source", "utmSource", "foo", " ", "bogus"}
	for _, in := range invalid {
		in := in
		t.Run(fmt.Sprintf("invalid_%q", in), func(t *testing.T) {
			t.Parallel()
			got, err := ValidateUtmType(in)
			if !errors.Is(err, ErrInvalidUtmType) {
				t.Errorf("got err=%v, want ErrInvalidUtmType", err)
			}
			if got != "" {
				t.Errorf("got type=%q, want empty", got)
			}
		})
	}
}

func TestTopUTM_DisabledService(t *testing.T) {
	t.Parallel()

	s, err := NewAnalyticsService("", "", "", "dev", 0, silentLogger())
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	if !s.Disabled() {
		t.Fatalf("expected disabled service")
	}
	if _, err := s.TopUTM(context.Background(), "7d", UtmSource, 10); !errors.Is(err, ErrServiceDisabled) {
		t.Errorf("got %v, want ErrServiceDisabled", err)
	}
}

func TestTopUTM_InputValidation(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Errorf("upstream hit during TopUTM validation test")
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	ctx := context.Background()

	t.Run("invalid range", func(t *testing.T) {
		t.Parallel()
		for _, r := range []string{"", "1d", "90d", "bogus"} {
			if _, err := s.TopUTM(ctx, r, UtmSource, 10); !errors.Is(err, ErrInvalidRange) {
				t.Errorf("range=%q: got %v, want ErrInvalidRange", r, err)
			}
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		t.Parallel()
		for _, l := range []int{0, -1, 101, 99999} {
			if _, err := s.TopUTM(ctx, "7d", UtmSource, l); !errors.Is(err, ErrInvalidLimit) {
				t.Errorf("limit=%d: got %v, want ErrInvalidLimit", l, err)
			}
		}
	})

	t.Run("zero utm type rejected", func(t *testing.T) {
		t.Parallel()
		if _, err := s.TopUTM(ctx, "7d", UtmType(""), 10); !errors.Is(err, ErrInvalidUtmType) {
			t.Errorf("zero UtmType: got %v, want ErrInvalidUtmType", err)
		}
		if _, err := s.TopUTM(ctx, "7d", UtmType("utm_bogus"), 10); !errors.Is(err, ErrInvalidUtmType) {
			t.Errorf("bogus UtmType: got %v, want ErrInvalidUtmType", err)
		}
	})
}

func TestTopUTM_HappyPath(t *testing.T) {
	t.Parallel()

	var gotType, gotLimit string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotType = r.URL.Query().Get("type")
		gotLimit = r.URL.Query().Get("limit")
		_, _ = io.WriteString(w, `[
          {"x": "instagram", "y": 420},
          {"x": "",          "y": 50},
          {"x": "  ",        "y": 7}
        ]`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	got, err := s.TopUTM(context.Background(), "30d", UtmSource, 10)
	if err != nil {
		t.Fatalf("TopUTM: %v", err)
	}

	if gotType != "utmSource" {
		t.Errorf("upstream type: got %q, want utmSource", gotType)
	}
	if gotLimit != "10" {
		t.Errorf("upstream limit: got %q, want 10", gotLimit)
	}

	if len(got) != 3 {
		t.Fatalf("len: got %d, want 3", len(got))
	}
	if got[0].Value != "instagram" || got[0].Views != 420 {
		t.Errorf("[0]: got %+v, want {instagram 420}", got[0])
	}
	if got[1].Value != "(none)" {
		t.Errorf("empty→(none): got %q", got[1].Value)
	}
	if got[1].Views != 50 {
		t.Errorf("[1].Views: got %d, want 50", got[1].Views)
	}
	if got[2].Value != "(none)" {
		t.Errorf("whitespace→(none): got %q", got[2].Value)
	}
}

func TestTopUTM_CacheByType(t *testing.T) {
	t.Parallel()

	var hits int32
	byType := map[string]int{}
	var mu sync.Mutex

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		mu.Lock()
		byType[r.URL.Query().Get("type")]++
		mu.Unlock()
		_, _ = io.WriteString(w, `[{"x":"a","y":1}]`)
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if _, err := s.TopUTM(ctx, "7d", UtmSource, 10); err != nil {
			t.Fatalf("source call %d: %v", i, err)
		}
	}
	for i := 0; i < 2; i++ {
		if _, err := s.TopUTM(ctx, "7d", UtmMedium, 10); err != nil {
			t.Fatalf("medium call %d: %v", i, err)
		}
	}
	if _, err := s.TopUTM(ctx, "7d", UtmCampaign, 10); err != nil {
		t.Fatalf("campaign call: %v", err)
	}

	if got := atomic.LoadInt32(&hits); got != 3 {
		t.Errorf("upstream hits: got %d, want 3 (one per distinct type)", got)
	}

	mu.Lock()
	defer mu.Unlock()
	for _, want := range []string{"utmSource", "utmMedium", "utmCampaign"} {
		if byType[want] != 1 {
			t.Errorf("upstream call count for %q: got %d, want 1", want, byType[want])
		}
	}
}

func TestGet_BodySizeCapped(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"pageviews":{"value":1,"prev":0},"visitors":{"value":1,"prev":0},"visits":{"value":1,"prev":0},"bounces":{"value":0,"prev":0},"totaltime":{"value":0,"prev":0},"filler":"`)
		chunk := strings.Repeat("A", 4096)
		const target = 2 * (1 << 20)
		written := 0
		for written < target {
			n, err := io.WriteString(w, chunk)
			if err != nil {
				return
			}
			written += n
		}
	}))
	t.Cleanup(ts.Close)

	s := newTestService(t, ts, time.Minute)

	done := make(chan error, 1)
	go func() {
		_, err := s.Summary(context.Background(), "7d")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected decode error on oversized body, got nil")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Summary did not return within 5s — body cap did not kick in")
	}
}
