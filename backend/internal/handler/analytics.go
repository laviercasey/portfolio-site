package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/portfolio/backend/internal/service"
)

const (
	defaultRange          = "30d"
	defaultLimit          = 10
	defaultTimeseriesUnit = "pageviews"

	errMsgNotConfigured   = "analytics_not_configured"
	errMsgUpstreamError   = "analytics upstream error"
	errMsgUnavailable     = "analytics unavailable"
	errMsgInternal        = "internal error"
	errMsgInvalidRange    = "invalid range"
	errMsgInvalidLimit    = "invalid limit"
	errMsgInvalidMetric   = "invalid metric"
	errMsgInvalidUtmType  = "invalid utm type"
	errMsgUtmTypeRequired = "utm type required"
)

var allowedPublicUtmTypes = map[string]struct{}{
	"source":   {},
	"medium":   {},
	"campaign": {},
}

type AnalyticsHandler struct {
	svc    *service.AnalyticsService
	logger *slog.Logger
}

func NewAnalyticsHandler(svc *service.AnalyticsService, logger *slog.Logger) *AnalyticsHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &AnalyticsHandler{svc: svc, logger: logger}
}

func parseRange(r *http.Request) string {
	v := r.URL.Query().Get("range")
	if v == "" {
		return defaultRange
	}
	return v
}

func parseLimit(r *http.Request) (int, bool) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return defaultLimit, true
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return n, true
}

func parseMetric(r *http.Request) string {
	v := r.URL.Query().Get("metric")
	if v == "" {
		return defaultTimeseriesUnit
	}
	return v
}

func (h *AnalyticsHandler) mapServiceError(
	w http.ResponseWriter,
	endpoint string,
	err error,
) {
	switch {
	case errors.Is(err, service.ErrServiceDisabled):
		writeAPIError(w, http.StatusServiceUnavailable, errMsgNotConfigured)
		return
	case errors.Is(err, service.ErrInvalidRange):
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidRange)
		return
	case errors.Is(err, service.ErrInvalidLimit):
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidLimit)
		return
	case errors.Is(err, service.ErrInvalidMetric):
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidMetric)
		return
	case errors.Is(err, service.ErrInvalidUtmType):
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidUtmType)
		return
	}

	var up *service.UpstreamError
	if errors.As(err, &up) {
		h.logger.Error(
			"analytics upstream failure",
			slog.String("endpoint", endpoint),
			slog.Int("status", up.Status),
			slog.String("error", err.Error()),
		)
		switch {
		case up.Status == 0,
			up.Status >= 500,
			up.Status == http.StatusUnauthorized,
			up.Status == http.StatusForbidden:
			writeAPIError(w, http.StatusServiceUnavailable, errMsgUnavailable)
		default:
			writeAPIError(w, http.StatusBadGateway, errMsgUpstreamError)
		}
		return
	}

	h.logger.Error(
		"analytics internal error",
		slog.String("endpoint", endpoint),
		slog.String("error", err.Error()),
	)
	writeAPIError(w, http.StatusInternalServerError, errMsgInternal)
}

func (h *AnalyticsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	rng := parseRange(r)

	out, err := h.svc.Summary(r.Context(), rng)
	if err != nil {
		h.mapServiceError(w, "summary", err)
		return
	}
	writeSuccess(w, out)
}

func (h *AnalyticsHandler) TopPages(w http.ResponseWriter, r *http.Request) {
	rng := parseRange(r)
	limit, ok := parseLimit(r)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidLimit)
		return
	}

	out, err := h.svc.TopPages(r.Context(), rng, limit)
	if err != nil {
		h.mapServiceError(w, "top-pages", err)
		return
	}
	writeSuccess(w, out)
}

func (h *AnalyticsHandler) TopReferrers(w http.ResponseWriter, r *http.Request) {
	rng := parseRange(r)
	limit, ok := parseLimit(r)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidLimit)
		return
	}

	out, err := h.svc.TopReferrers(r.Context(), rng, limit)
	if err != nil {
		h.mapServiceError(w, "top-referrers", err)
		return
	}
	writeSuccess(w, out)
}

func (h *AnalyticsHandler) TopCountries(w http.ResponseWriter, r *http.Request) {
	rng := parseRange(r)
	limit, ok := parseLimit(r)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidLimit)
		return
	}

	out, err := h.svc.TopCountries(r.Context(), rng, limit)
	if err != nil {
		h.mapServiceError(w, "top-countries", err)
		return
	}
	writeSuccess(w, out)
}

func (h *AnalyticsHandler) TopUTM(w http.ResponseWriter, r *http.Request) {
	rng := parseRange(r)
	limit, ok := parseLimit(r)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidLimit)
		return
	}

	typeStr := r.URL.Query().Get("type")
	if typeStr == "" {
		writeAPIError(w, http.StatusBadRequest, errMsgUtmTypeRequired)
		return
	}
	if _, allowed := allowedPublicUtmTypes[typeStr]; !allowed {
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidUtmType)
		return
	}
	utmType, err := service.ValidateUtmType(typeStr)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, errMsgInvalidUtmType)
		return
	}

	out, err := h.svc.TopUTM(r.Context(), rng, utmType, limit)
	if err != nil {
		h.mapServiceError(w, "top-utm", err)
		return
	}
	writeSuccess(w, out)
}

func (h *AnalyticsHandler) Timeseries(w http.ResponseWriter, r *http.Request) {
	rng := parseRange(r)
	metric := parseMetric(r)

	out, err := h.svc.Timeseries(r.Context(), rng, metric)
	if err != nil {
		h.mapServiceError(w, "timeseries", err)
		return
	}
	writeSuccess(w, out)
}
