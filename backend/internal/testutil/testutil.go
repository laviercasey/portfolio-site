package testutil

import (
	"io"
	"log/slog"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/portfolio/backend/internal/database"
)

func SilentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
}

func NewMockPool(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func NewMockDB(t *testing.T) (*database.DB, pgxmock.PgxPoolIface) {
	t.Helper()
	pool := NewMockPool(t)
	return database.NewWithPool(pool), pool
}

func AnyArgs(n int) []any {
	out := make([]any, n)
	for i := range out {
		out[i] = pgxmock.AnyArg()
	}
	return out
}
