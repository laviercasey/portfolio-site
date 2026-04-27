package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestWriteJSON_Success(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	writeJSON(rec, 201, map[string]string{"hello": "world"})

	if rec.Code != 201 {
		t.Fatalf("status: got %d, want 201", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("content-type: %q", ct)
	}
	var out map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["hello"] != "world" {
		t.Errorf("body: got %v", out)
	}
}

func TestWriteJSON_MarshalError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	ch := make(chan int)
	writeJSON(rec, 200, ch)

	if rec.Code != 500 {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "internal") {
		t.Errorf("expected internal error, got %q", rec.Body.String())
	}
}

func TestWriteError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	writeError(rec, 400, "bad thing")

	if rec.Code != 400 {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	var out map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["error"] != "bad thing" {
		t.Errorf("error: got %q", out["error"])
	}
}

type sampleInput struct {
	Name  string `validate:"required,min=3"`
	Email string `validate:"required,email"`
}

func TestFormatValidationError_MultipleFields(t *testing.T) {
	t.Parallel()

	v := validator.New()
	err := v.Struct(sampleInput{Name: "a", Email: "not-email"})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	msg := formatValidationError(err)
	if !strings.Contains(msg, "Name") {
		t.Errorf("expected Name in msg: %q", msg)
	}
	if !strings.Contains(msg, "Email") {
		t.Errorf("expected Email in msg: %q", msg)
	}
	if !strings.Contains(msg, ";") {
		t.Errorf("expected semicolon separator: %q", msg)
	}
}

func TestFormatValidationError_NonValidationError(t *testing.T) {
	t.Parallel()

	msg := formatValidationError(errors.New("random"))
	if msg != "validation failed" {
		t.Errorf("got %q, want 'validation failed'", msg)
	}
}
