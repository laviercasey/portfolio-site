package handler

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type recordingRevalidator struct {
	mu    sync.Mutex
	calls [][]string
	err   error
	done  chan struct{}
}

func newRecordingRevalidator() *recordingRevalidator {
	return &recordingRevalidator{done: make(chan struct{}, 4)}
}

func (r *recordingRevalidator) Revalidate(_ context.Context, paths []string) error {
	r.mu.Lock()
	r.calls = append(r.calls, paths)
	r.mu.Unlock()
	r.done <- struct{}{}
	return r.err
}

func (r *recordingRevalidator) Calls() [][]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([][]string, len(r.calls))
	copy(out, r.calls)
	return out
}

func TestRevalidateAsync_NilNoop(t *testing.T) {
	t.Parallel()
	revalidateAsync(nil, []string{"/"})
}

func TestRevalidateAsync_FiresInBackground(t *testing.T) {
	t.Parallel()
	rec := newRecordingRevalidator()
	revalidateAsync(rec, []string{"/", "/x"})

	select {
	case <-rec.done:
	case <-time.After(2 * time.Second):
		t.Fatal("revalidate not called in time")
	}

	calls := rec.Calls()
	if len(calls) != 1 {
		t.Fatalf("calls = %d, want 1", len(calls))
	}
	if len(calls[0]) != 2 || calls[0][0] != "/" || calls[0][1] != "/x" {
		t.Errorf("paths = %v", calls[0])
	}
}

func TestRevalidateAsync_ErrorSwallowed(t *testing.T) {
	t.Parallel()
	rec := newRecordingRevalidator()
	rec.err = errors.New("boom")
	revalidateAsync(rec, []string{"/"})

	select {
	case <-rec.done:
	case <-time.After(2 * time.Second):
		t.Fatal("not called")
	}
}
