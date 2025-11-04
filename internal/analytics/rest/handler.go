package rest

import (
	"context"
	"fmt"
	"github.com/ilam072/sales-tracker/internal/response"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
	"strconv"
	"time"
)

type Analytics interface {
	Sum(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error)
	Avg(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error)
	Count(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (int, error)
	Median(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error)
	PercentileNinetieth(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error)
}

type Validator interface {
	Validate(i interface{}) error
}

type AnalyticsHandler struct {
	analytics Analytics
	validator Validator
}

func NewAnalyticsHandler(analytics Analytics, validator Validator) *AnalyticsHandler {
	return &AnalyticsHandler{analytics: analytics, validator: validator}
}

func (h *AnalyticsHandler) Sum(c *ginext.Context) {
	from, to, categoryID, itemType, err := parseQueryParams(c)
	if err != nil {
		response.Error(err.Error()).WriteJSON(c, http.StatusBadRequest)
		return
	}

	sum, err := h.analytics.Sum(c.Request.Context(), from, to, categoryID, itemType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate sum")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"sum": sum})
}

func (h *AnalyticsHandler) Avg(c *ginext.Context) {
	from, to, categoryID, itemType, err := parseQueryParams(c)
	if err != nil {
		response.Error(err.Error()).WriteJSON(c, http.StatusBadRequest)
		return
	}

	avg, err := h.analytics.Avg(c.Request.Context(), from, to, categoryID, itemType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate average")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"average": avg})
}

func (h *AnalyticsHandler) Count(c *ginext.Context) {
	from, to, categoryID, itemType, err := parseQueryParams(c)
	if err != nil {
		response.Error(err.Error()).WriteJSON(c, http.StatusBadRequest)
		return
	}

	count, err := h.analytics.Count(c.Request.Context(), from, to, categoryID, itemType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to count items")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"count": count})
}

func (h *AnalyticsHandler) Median(c *ginext.Context) {
	from, to, categoryID, itemType, err := parseQueryParams(c)
	if err != nil {
		response.Error(err.Error()).WriteJSON(c, http.StatusBadRequest)
		return
	}

	median, err := h.analytics.Median(c.Request.Context(), from, to, categoryID, itemType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate median")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"median": median})
}

func (h *AnalyticsHandler) PercentileNinetieth(c *ginext.Context) {
	from, to, categoryID, itemType, err := parseQueryParams(c)
	if err != nil {
		response.Error(err.Error()).WriteJSON(c, http.StatusBadRequest)
		return
	}

	p90, err := h.analytics.PercentileNinetieth(c.Request.Context(), from, to, categoryID, itemType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate 90th percentile")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"percentile_90": p90})
}

// parseQueryParams парсит query параметры ?from=...&to=...&category_id=...&type=...
func parseQueryParams(c *ginext.Context) (from, to *time.Time, categoryID *int, itemType *string, err error) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	categoryStr := c.Query("category_id")
	typeStr := c.Query("type")

	if fromStr != "" {
		t, e := time.Parse(time.DateOnly, fromStr)
		if e != nil {
			return nil, nil, nil, nil, fmt.Errorf("invalid 'from' format, expected YYYY-MM-DD")
		}
		from = &t
	}

	if toStr != "" {
		t, e := time.Parse(time.DateOnly, toStr)
		if e != nil {
			return nil, nil, nil, nil, fmt.Errorf("invalid 'to' format, expected YYYY-MM-DD")
		}
		to = &t
	}

	if categoryStr != "" {
		id, e := strconv.Atoi(categoryStr)
		if e != nil {
			return nil, nil, nil, nil, fmt.Errorf("invalid 'category_id', must be integer")
		}
		categoryID = &id
	}

	if typeStr != "" {
		itemType = &typeStr
	}

	return from, to, categoryID, itemType, nil
}
