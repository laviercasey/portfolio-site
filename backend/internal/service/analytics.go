package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/portfolio/backend/internal/model"
)

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var (
	ErrServiceDisabled = errors.New("analytics: service disabled")
	ErrInvalidRange    = errors.New("analytics: invalid range")
	ErrInvalidLimit    = errors.New("analytics: invalid limit")
	ErrInvalidMetric   = errors.New("analytics: invalid metric")
	ErrInvalidUtmType  = errors.New("analytics: invalid utm type")
	ErrUpstream        = errors.New("analytics: upstream error")
)

type UtmType string

const (
	UtmSource   UtmType = "utmSource"
	UtmMedium   UtmType = "utmMedium"
	UtmCampaign UtmType = "utmCampaign"
	UtmTerm     UtmType = "utmTerm"
	UtmContent  UtmType = "utmContent"
)

func ValidateUtmType(t string) (UtmType, error) {
	switch t {
	case "source":
		return UtmSource, nil
	case "medium":
		return UtmMedium, nil
	case "campaign":
		return UtmCampaign, nil
	case "term":
		return UtmTerm, nil
	case "content":
		return UtmContent, nil
	default:
		return "", ErrInvalidUtmType
	}
}

type UpstreamError struct {
	Status int
	Err    error
}

func (e *UpstreamError) Error() string {
	if e.Status == 0 {
		return fmt.Sprintf("analytics upstream: %v", e.Err)
	}
	return fmt.Sprintf("analytics upstream: status=%d: %v", e.Status, e.Err)
}

func (e *UpstreamError) Unwrap() error { return e.Err }

func (e *UpstreamError) Is(target error) bool {
	return target == ErrUpstream
}

const (
	umamiClientTimeout = 5 * time.Second
	defaultCacheTTL    = 60 * time.Second
	maxLimit           = 100
)

type AnalyticsService struct {
	client    *http.Client
	apiURL    *url.URL
	apiKey    string
	websiteID string
	cache     *ttlCache
	ttl       time.Duration
	logger    *slog.Logger
	disabled  bool
}

func NewAnalyticsService(
	apiURL, apiKey, websiteID, env string,
	cacheTTL time.Duration,
	logger *slog.Logger,
) (*AnalyticsService, error) {
	if logger == nil {
		logger = slog.Default()
	}
	if apiURL == "" || apiKey == "" || websiteID == "" {
		logger.Info("analytics service: disabled (missing envs)")
		return &AnalyticsService{
			disabled: true,
			logger:   logger,
		}, nil
	}

	if !uuidRe.MatchString(websiteID) {
		return nil, fmt.Errorf("analytics: UMAMI_WEBSITE_ID must be a valid UUID")
	}

	u, err := validateUmamiURL(apiURL, env)
	if err != nil {
		return nil, fmt.Errorf("analytics: validate url: %w", err)
	}

	ttl := cacheTTL
	if ttl <= 0 {
		ttl = defaultCacheTTL
	}

	client := &http.Client{
		Timeout: umamiClientTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	logger.Info(
		"analytics service: enabled",
		slog.String("url", u.Scheme+"://"+u.Host),
	)

	return &AnalyticsService{
		client:    client,
		apiURL:    u,
		apiKey:    apiKey,
		websiteID: websiteID,
		cache:     newTTLCache(),
		ttl:       ttl,
		logger:    logger,
	}, nil
}

func (s *AnalyticsService) Disabled() bool { return s.disabled }

func validateRange(r string) error {
	switch r {
	case "7d", "30d":
		return nil
	default:
		return ErrInvalidRange
	}
}

func validateLimit(l int) error {
	if l < 1 || l > maxLimit {
		return ErrInvalidLimit
	}
	return nil
}

func validateMetric(m string) error {
	switch m {
	case "pageviews", "visitors":
		return nil
	default:
		return ErrInvalidMetric
	}
}

func rangeWindow(r string) (startMs, endMs int64) {
	now := time.Now().UTC()
	var d time.Duration
	switch r {
	case "7d":
		d = 7 * 24 * time.Hour
	case "30d":
		d = 30 * 24 * time.Hour
	}
	return now.Add(-d).UnixMilli(), now.UnixMilli()
}

func (s *AnalyticsService) Summary(
	ctx context.Context,
	rng string,
) (*model.AnalyticsSummary, error) {
	if s.disabled {
		return nil, ErrServiceDisabled
	}
	if err := validateRange(rng); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("summary:%s:0", rng)
	out, err := s.cached(ctx, key, func(ctx context.Context) (any, error) {
		return s.fetchSummary(ctx, rng)
	})
	if err != nil {
		return nil, fmt.Errorf("analytics: summary: %w", err)
	}
	typed, ok := out.(*model.AnalyticsSummary)
	if !ok {
		return nil, fmt.Errorf("analytics: summary: unexpected cache value type %T", out)
	}
	return typed, nil
}

func (s *AnalyticsService) TopPages(
	ctx context.Context,
	rng string,
	limit int,
) ([]model.TopPage, error) {
	if s.disabled {
		return nil, ErrServiceDisabled
	}
	if err := validateRange(rng); err != nil {
		return nil, err
	}
	if err := validateLimit(limit); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("top-pages:%s:%d", rng, limit)
	out, err := s.cached(ctx, key, func(ctx context.Context) (any, error) {
		return s.fetchTopPages(ctx, rng, limit)
	})
	if err != nil {
		return nil, fmt.Errorf("analytics: top-pages: %w", err)
	}
	typed, ok := out.([]model.TopPage)
	if !ok {
		return nil, fmt.Errorf("analytics: top-pages: unexpected cache value type %T", out)
	}
	return typed, nil
}

func (s *AnalyticsService) TopReferrers(
	ctx context.Context,
	rng string,
	limit int,
) ([]model.TopReferrer, error) {
	if s.disabled {
		return nil, ErrServiceDisabled
	}
	if err := validateRange(rng); err != nil {
		return nil, err
	}
	if err := validateLimit(limit); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("top-referrers:%s:%d", rng, limit)
	out, err := s.cached(ctx, key, func(ctx context.Context) (any, error) {
		return s.fetchTopReferrers(ctx, rng, limit)
	})
	if err != nil {
		return nil, fmt.Errorf("analytics: top-referrers: %w", err)
	}
	typed, ok := out.([]model.TopReferrer)
	if !ok {
		return nil, fmt.Errorf("analytics: top-referrers: unexpected cache value type %T", out)
	}
	return typed, nil
}

func (s *AnalyticsService) TopCountries(
	ctx context.Context,
	rng string,
	limit int,
) ([]model.TopCountry, error) {
	if s.disabled {
		return nil, ErrServiceDisabled
	}
	if err := validateRange(rng); err != nil {
		return nil, err
	}
	if err := validateLimit(limit); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("top-countries:%s:%d", rng, limit)
	out, err := s.cached(ctx, key, func(ctx context.Context) (any, error) {
		return s.fetchTopCountries(ctx, rng, limit)
	})
	if err != nil {
		return nil, fmt.Errorf("analytics: top-countries: %w", err)
	}
	typed, ok := out.([]model.TopCountry)
	if !ok {
		return nil, fmt.Errorf("analytics: top-countries: unexpected cache value type %T", out)
	}
	return typed, nil
}

func (s *AnalyticsService) TopUTM(
	ctx context.Context,
	rng string,
	utmType UtmType,
	limit int,
) ([]model.TopUtm, error) {
	if s.disabled {
		return nil, ErrServiceDisabled
	}
	if err := validateRange(rng); err != nil {
		return nil, err
	}
	if err := validateLimit(limit); err != nil {
		return nil, err
	}
	switch utmType {
	case UtmSource, UtmMedium, UtmCampaign, UtmTerm, UtmContent:
	default:
		return nil, ErrInvalidUtmType
	}

	key := fmt.Sprintf("top-utm:%s:%s:%d", rng, utmType, limit)
	out, err := s.cached(ctx, key, func(ctx context.Context) (any, error) {
		return s.fetchTopUTM(ctx, rng, utmType, limit)
	})
	if err != nil {
		return nil, fmt.Errorf("analytics: top-utm: %w", err)
	}
	typed, ok := out.([]model.TopUtm)
	if !ok {
		return nil, fmt.Errorf("analytics: top-utm: unexpected cache value type %T", out)
	}
	return typed, nil
}

func (s *AnalyticsService) Timeseries(
	ctx context.Context,
	rng, metric string,
) ([]model.TimeseriesPoint, error) {
	if s.disabled {
		return nil, ErrServiceDisabled
	}
	if err := validateRange(rng); err != nil {
		return nil, err
	}
	if err := validateMetric(metric); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("timeseries:%s:%s", rng, metric)
	out, err := s.cached(ctx, key, func(ctx context.Context) (any, error) {
		return s.fetchTimeseries(ctx, rng, metric)
	})
	if err != nil {
		return nil, fmt.Errorf("analytics: timeseries: %w", err)
	}
	typed, ok := out.([]model.TimeseriesPoint)
	if !ok {
		return nil, fmt.Errorf("analytics: timeseries: unexpected cache value type %T", out)
	}
	return typed, nil
}

type umamiStatComparison struct {
	Pageviews int64 `json:"pageviews"`
	Visitors  int64 `json:"visitors"`
	Visits    int64 `json:"visits"`
	Bounces   int64 `json:"bounces"`
	TotalTime int64 `json:"totaltime"`
}

type umamiStats struct {
	Pageviews  int64               `json:"pageviews"`
	Visitors   int64               `json:"visitors"`
	Visits     int64               `json:"visits"`
	Bounces    int64               `json:"bounces"`
	TotalTime  int64               `json:"totaltime"`
	Comparison umamiStatComparison `json:"comparison"`
}

type umamiMetric struct {
	X string      `json:"x"`
	Y json.Number `json:"y"`
}

type umamiPageviews struct {
	Pageviews []umamiMetric `json:"pageviews"`
	Sessions  []umamiMetric `json:"sessions"`
}

func (s *AnalyticsService) fetchSummary(
	ctx context.Context,
	rng string,
) (*model.AnalyticsSummary, error) {
	startMs, endMs := rangeWindow(rng)

	q := url.Values{}
	q.Set("startAt", strconv.FormatInt(startMs, 10))
	q.Set("endAt", strconv.FormatInt(endMs, 10))
	q.Set("unit", "day")
	q.Set("timezone", "UTC")
	q.Set("compare", "prev")

	var stats umamiStats
	if err := s.get(ctx, s.websitePath("stats"), q, &stats); err != nil {
		return nil, err
	}

	bounceRate := 0.0
	if stats.Visits > 0 {
		bounceRate = float64(stats.Bounces) / float64(stats.Visits)
	}
	if bounceRate < 0 {
		bounceRate = 0
	}
	if bounceRate > 1 {
		bounceRate = 1
	}

	avgSession := int64(0)
	if stats.Visits > 0 {
		avgSession = stats.TotalTime / stats.Visits
	}
	if avgSession < 0 {
		avgSession = 0
	}

	prev := stats.Comparison.Pageviews
	delta := 0.0
	switch {
	case prev == 0 && stats.Pageviews == 0:
		delta = 0
	case prev == 0:
		delta = 10
	default:
		delta = (float64(stats.Pageviews) - float64(prev)) / float64(prev)
	}
	if delta > 10 {
		delta = 10
	}
	if delta < -10 {
		delta = -10
	}
	if math.IsNaN(delta) || math.IsInf(delta, 0) {
		delta = 0
	}

	return &model.AnalyticsSummary{
		Range:             rng,
		Pageviews:         stats.Pageviews,
		UniqueVisitors:    stats.Visitors,
		BounceRate:        bounceRate,
		AvgSessionSeconds: avgSession,
		DeltaPageviews:    delta,
		PreviousPageviews: prev,
	}, nil
}

func (s *AnalyticsService) fetchTopPages(
	ctx context.Context,
	rng string,
	limit int,
) ([]model.TopPage, error) {
	metrics, err := s.fetchMetrics(ctx, rng, "path", limit)
	if err != nil {
		return nil, err
	}

	out := make([]model.TopPage, 0, len(metrics))
	for _, m := range metrics {
		path := normalizePath(m.X)
		views, _ := m.Y.Int64()
		out = append(out, model.TopPage{
			Path:    path,
			Views:   views,
			Uniques: 0,
		})
	}
	return out, nil
}

func (s *AnalyticsService) fetchTopReferrers(
	ctx context.Context,
	rng string,
	limit int,
) ([]model.TopReferrer, error) {
	metrics, err := s.fetchMetrics(ctx, rng, "referrer", limit)
	if err != nil {
		return nil, err
	}

	out := make([]model.TopReferrer, 0, len(metrics))
	for _, m := range metrics {
		ref := strings.TrimSpace(m.X)
		if ref == "" {
			ref = "(direct)"
		}
		views, _ := m.Y.Int64()
		out = append(out, model.TopReferrer{
			Referrer: ref,
			Views:    views,
		})
	}
	return out, nil
}

func (s *AnalyticsService) fetchTopCountries(
	ctx context.Context,
	rng string,
	limit int,
) ([]model.TopCountry, error) {
	metrics, err := s.fetchMetrics(ctx, rng, "country", limit)
	if err != nil {
		return nil, err
	}

	out := make([]model.TopCountry, 0, len(metrics))
	for _, m := range metrics {
		c := strings.ToUpper(strings.TrimSpace(m.X))
		if len(c) != 2 {
			c = "ZZ"
		}
		views, _ := m.Y.Int64()
		out = append(out, model.TopCountry{
			Country: c,
			Views:   views,
		})
	}
	return out, nil
}

func (s *AnalyticsService) fetchTopUTM(
	ctx context.Context,
	rng string,
	utmType UtmType,
	limit int,
) ([]model.TopUtm, error) {
	metrics, err := s.fetchMetrics(ctx, rng, string(utmType), limit)
	if err != nil {
		return nil, err
	}

	out := make([]model.TopUtm, 0, len(metrics))
	for _, m := range metrics {
		val := strings.TrimSpace(m.X)
		if val == "" {
			val = "(none)"
		}
		views, _ := m.Y.Int64()
		out = append(out, model.TopUtm{
			Value: val,
			Views: views,
		})
	}
	return out, nil
}

func (s *AnalyticsService) fetchMetrics(
	ctx context.Context,
	rng, kind string,
	limit int,
) ([]umamiMetric, error) {
	startMs, endMs := rangeWindow(rng)
	q := url.Values{}
	q.Set("startAt", strconv.FormatInt(startMs, 10))
	q.Set("endAt", strconv.FormatInt(endMs, 10))
	q.Set("unit", "day")
	q.Set("timezone", "UTC")
	q.Set("type", kind)
	q.Set("limit", strconv.Itoa(limit))

	var out []umamiMetric
	if err := s.get(ctx, s.websitePath("metrics"), q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *AnalyticsService) fetchTimeseries(
	ctx context.Context,
	rng, metric string,
) ([]model.TimeseriesPoint, error) {
	startMs, endMs := rangeWindow(rng)
	q := url.Values{}
	q.Set("startAt", strconv.FormatInt(startMs, 10))
	q.Set("endAt", strconv.FormatInt(endMs, 10))
	q.Set("unit", "day")
	q.Set("timezone", "UTC")

	var payload umamiPageviews
	if err := s.get(ctx, s.websitePath("pageviews"), q, &payload); err != nil {
		return nil, err
	}

	series := payload.Pageviews
	if metric == "visitors" {
		series = payload.Sessions
	}

	out := make([]model.TimeseriesPoint, 0, len(series))
	for _, p := range series {
		date := normalizeDate(p.X)
		v, _ := p.Y.Int64()
		if v < 0 {
			v = 0
		}
		out = append(out, model.TimeseriesPoint{Date: date, Value: v})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Date < out[j].Date })
	return out, nil
}

func (s *AnalyticsService) websitePath(sub string) string {
	return "/api/websites/" + s.websiteID + "/" + sub
}

func (s *AnalyticsService) get(
	ctx context.Context,
	path string,
	q url.Values,
	dst any,
) error {
	rel := &url.URL{Path: path, RawQuery: q.Encode()}
	full := s.apiURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full.String(), nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return &UpstreamError{Status: 0, Err: err}
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		s.logger.Warn(
			"umami non-2xx",
			slog.String("path", path),
			slog.Int("status", resp.StatusCode),
			slog.String("body_prefix", truncateForLog(string(buf))),
		)
		return &UpstreamError{
			Status: resp.StatusCode,
			Err:    fmt.Errorf("status %d", resp.StatusCode),
		}
	}

	const maxResponseBytes = 1 << 20
	dec := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes))
	dec.UseNumber()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode upstream: %w", err)
	}
	return nil
}

func truncateForLog(s string) string {
	const max = 200
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

func normalizePath(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "/"
	}
	u, err := url.Parse(s)
	if err != nil || u.Path == "" {
		if i := strings.IndexAny(s, "?#"); i >= 0 {
			s = s[:i]
		}
		if s == "" {
			return "/"
		}
		return s
	}
	if u.Path == "" {
		return "/"
	}
	return u.Path
}

func normalizeDate(raw string) string {
	s := strings.TrimSpace(raw)
	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC().Format("2006-01-02")
		}
	}
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

type cacheEntry struct {
	value     any
	expiresAt time.Time
}

type ttlCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry

	flightMu sync.Mutex
	flights  map[string]*inflight
}

type inflight struct {
	done  chan struct{}
	value any
	err   error
}

func newTTLCache() *ttlCache {
	return &ttlCache{
		entries: make(map[string]cacheEntry),
		flights: make(map[string]*inflight),
	}
}

func (c *ttlCache) get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *ttlCache) set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}

func (s *AnalyticsService) cached(
	ctx context.Context,
	key string,
	fetch func(ctx context.Context) (any, error),
) (any, error) {
	if v, ok := s.cache.get(key); ok {
		return v, nil
	}

	s.cache.flightMu.Lock()
	if f, ok := s.cache.flights[key]; ok {
		s.cache.flightMu.Unlock()
		select {
		case <-f.done:
			return f.value, f.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	f := &inflight{done: make(chan struct{})}
	s.cache.flights[key] = f
	s.cache.flightMu.Unlock()

	defer func() {
		s.cache.flightMu.Lock()
		delete(s.cache.flights, key)
		s.cache.flightMu.Unlock()
		close(f.done)
	}()

	v, err := fetch(ctx)
	if err != nil {
		f.err = err
		return nil, err
	}
	f.value = v
	s.cache.set(key, v, s.ttl)
	return v, nil
}
