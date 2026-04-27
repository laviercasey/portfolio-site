package service

import (
	"testing"
	"time"
)

func TestNullableStr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    string
		isNil bool
		want  string
	}{
		{"empty returns nil", "", true, ""},
		{"non-empty returns pointer", "hello", false, "hello"},
		{"whitespace is non-empty", " ", false, " "},
		{"unicode non-empty", "тест", false, "тест"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := nullableStr(tc.in)
			if tc.isNil {
				if got != nil {
					t.Fatalf("expected nil, got %q", *got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected non-nil pointer")
			}
			if *got != tc.want {
				t.Errorf("got %q, want %q", *got, tc.want)
			}
		})
	}
}

func TestParseFlexibleDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		wantErr bool
		wantY   int
		wantM   time.Month
		wantD   int
	}{
		{"rfc3339", "2025-03-15T12:34:56Z", false, 2025, time.March, 15},
		{"iso date z", "2025-03-15T12:34:56Z", false, 2025, time.March, 15},
		{"iso date", "2025-03-15", false, 2025, time.March, 15},
		{"year month", "2025-03", false, 2025, time.March, 1},
		{"year only", "2025", false, 2025, time.January, 1},
		{"empty", "", true, 0, 0, 0},
		{"garbage", "not-a-date", true, 0, 0, 0},
		{"bad month", "2025-13-01", true, 0, 0, 0},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseFlexibleDate(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (value=%v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Year() != tc.wantY {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantY)
			}
			if got.Month() != tc.wantM {
				t.Errorf("month: got %v, want %v", got.Month(), tc.wantM)
			}
			if got.Day() != tc.wantD {
				t.Errorf("day: got %d, want %d", got.Day(), tc.wantD)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("map round trip", func(t *testing.T) {
		t.Parallel()
		in := map[string]string{"ru": "Привет", "en": "Hello"}
		b, err := marshalJSON(in)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if len(b) == 0 {
			t.Fatal("empty bytes")
		}
	})

	t.Run("slice", func(t *testing.T) {
		t.Parallel()
		in := []string{"a", "b"}
		b, err := marshalJSON(in)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(b) != `["a","b"]` {
			t.Errorf("got %q, want %q", string(b), `["a","b"]`)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		b, err := marshalJSON(nil)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(b) != "null" {
			t.Errorf("got %q, want null", string(b))
		}
	})

	t.Run("unmarshalable value returns error", func(t *testing.T) {
		t.Parallel()
		_, err := marshalJSON(make(chan int))
		if err == nil {
			t.Fatal("expected error for channel type")
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("empty data leaves dst untouched", func(t *testing.T) {
		t.Parallel()
		var dst map[string]string
		if err := unmarshalJSON(nil, &dst); err != nil {
			t.Fatalf("err: %v", err)
		}
		if dst != nil {
			t.Errorf("dst mutated: %v", dst)
		}
	})

	t.Run("zero bytes is treated as empty", func(t *testing.T) {
		t.Parallel()
		var dst map[string]string
		if err := unmarshalJSON([]byte{}, &dst); err != nil {
			t.Fatalf("err: %v", err)
		}
		if dst != nil {
			t.Errorf("dst mutated: %v", dst)
		}
	})

	t.Run("valid json populates dst", func(t *testing.T) {
		t.Parallel()
		var dst map[string]string
		if err := unmarshalJSON([]byte(`{"ru":"Привет"}`), &dst); err != nil {
			t.Fatalf("err: %v", err)
		}
		if dst["ru"] != "Привет" {
			t.Errorf("got %q, want Привет", dst["ru"])
		}
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		t.Parallel()
		var dst map[string]string
		if err := unmarshalJSON([]byte(`{bad}`), &dst); err == nil {
			t.Fatal("expected error")
		}
	})
}
