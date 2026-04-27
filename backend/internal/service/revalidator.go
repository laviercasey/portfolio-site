package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const revalidateClientTimeout = 5 * time.Second

type RevalidateService struct {
	url    string
	secret string
	client *http.Client
	logger *slog.Logger
}

func NewRevalidateService(url, secret string, logger *slog.Logger) *RevalidateService {
	if logger == nil {
		logger = slog.Default()
	}
	return &RevalidateService{
		url:    url,
		secret: secret,
		client: &http.Client{Timeout: revalidateClientTimeout},
		logger: logger,
	}
}

func (c *RevalidateService) Enabled() bool {
	return c.url != "" && c.secret != ""
}

func (c *RevalidateService) Revalidate(ctx context.Context, paths []string) error {
	if !c.Enabled() {
		return nil
	}
	if len(paths) == 0 {
		return nil
	}

	payload, err := json.Marshal(map[string][]string{"paths": paths})
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	c.logger.Info("revalidate: ok", slog.Any("paths", paths))
	return nil
}
