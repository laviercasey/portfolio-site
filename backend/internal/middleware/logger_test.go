package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func captureLogger(buf *bytes.Buffer) *slog.Logger {
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(h)
}

func parseLogLine(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse log line: %v (raw=%q)", err, string(raw))
	}
	return m
}

func TestLogger_LogsBasicFieldsOnSuccess(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := captureLogger(&buf)

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/widgets", nil)
	req.RemoteAddr = "203.0.113.7:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}

	m := parseLogLine(t, buf.Bytes())

	if m["msg"] != "request" {
		t.Errorf("msg = %v; want request", m["msg"])
	}
	if m["method"] != http.MethodGet {
		t.Errorf("method = %v; want GET", m["method"])
	}
	if m["path"] != "/api/widgets" {
		t.Errorf("path = %v; want /api/widgets", m["path"])
	}
	if s, ok := m["status"].(float64); !ok || int(s) != 200 {
		t.Errorf("status = %v; want 200", m["status"])
	}
	if sz, ok := m["size"].(float64); !ok || int(sz) != len(`{"ok":true}`) {
		t.Errorf("size = %v; want %d", m["size"], len(`{"ok":true}`))
	}
	if _, ok := m["duration"]; !ok {
		t.Error("duration field missing")
	}
	if m["remote"] != "203.0.113.7" {
		t.Errorf("remote = %v; want 203.0.113.7", m["remote"])
	}
	if m["level"] != "INFO" {
		t.Errorf("level = %v; want INFO", m["level"])
	}
}

func TestLogger_CapturesCustomStatusCode(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := captureLogger(&buf)

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("nope"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/teapot", nil)
	req.RemoteAddr = "198.51.100.4:9000"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Fatalf("status = %d; want 418", rec.Code)
	}

	m := parseLogLine(t, buf.Bytes())
	if s, ok := m["status"].(float64); !ok || int(s) != http.StatusTeapot {
		t.Errorf("logged status = %v; want 418", m["status"])
	}
	if m["method"] != http.MethodPost {
		t.Errorf("method = %v; want POST", m["method"])
	}
	if sz, ok := m["size"].(float64); !ok || int(sz) != len("nope") {
		t.Errorf("size = %v; want %d", m["size"], len("nope"))
	}
}

func TestLogger_DefaultStatusIs200WhenNoWriteHeader(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := captureLogger(&buf)

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	req := httptest.NewRequest(http.MethodGet, "/noop", nil)
	req.RemoteAddr = "192.0.2.1:80"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	m := parseLogLine(t, buf.Bytes())
	if s, ok := m["status"].(float64); !ok || int(s) != http.StatusOK {
		t.Errorf("status = %v; want 200 default", m["status"])
	}
	if sz, ok := m["size"].(float64); !ok || int(sz) != 0 {
		t.Errorf("size = %v; want 0", m["size"])
	}
}

func TestLogger_WritesExactlyOneLinePerRequest(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := captureLogger(&buf)

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodDelete, "/drop", nil)
	req.RemoteAddr = "192.0.2.9:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	lines := bytes.Count(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	if lines != 0 {
		t.Errorf("log line count = %d+1; want exactly 1", lines)
	}
	if buf.Len() == 0 {
		t.Fatal("logger wrote nothing")
	}
}
