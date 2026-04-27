package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/portfolio/backend/internal/service"
)

const revalidateAsyncTimeout = 5 * time.Second

func revalidateAsync(rev service.Revalidator, paths []string) {
	if rev == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), revalidateAsyncTimeout)
		defer cancel()
		if err := rev.Revalidate(ctx, paths); err != nil {
			slog.Warn("handler revalidate failed", slog.String("error", err.Error()))
		}
	}()
}
