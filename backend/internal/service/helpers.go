package service

import (
	"encoding/json"
	"fmt"
	"time"
)

func nullableStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseFlexibleDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"2006-01",
		"2006",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date %q", s)
}

func marshalJSON(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal %T: %w", v, err)
	}
	return b, nil
}

func unmarshalJSON(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal %T: %w", v, err)
	}
	return nil
}
