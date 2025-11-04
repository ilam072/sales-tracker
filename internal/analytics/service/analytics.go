package service

import (
	"context"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/pkg/errutils"
	"time"
)

type AnalyticsRepo interface {
	Sum(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error)
	Avg(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error)
	Count(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (int, error)
	Median(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error)
	PercentileNinetieth(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error)
}

type Analytics struct {
	repo AnalyticsRepo
}

func New(repo AnalyticsRepo) *Analytics {
	return &Analytics{repo: repo}
}

func (a *Analytics) Sum(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error) {
	const op = "service.analytics.Sum"

	var typeFilter *domain.ItemType
	if itemType != nil {
		t := domain.ItemType(*itemType)
		typeFilter = &t
	}

	sum, err := a.repo.Sum(ctx, from, to, categoryID, typeFilter)
	if err != nil {
		return 0, errutils.Wrap(op, err)
	}

	return sum, nil
}

func (a *Analytics) Avg(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error) {
	const op = "service.analytics.Avg"

	var typeFilter *domain.ItemType
	if itemType != nil {
		t := domain.ItemType(*itemType)
		typeFilter = &t
	}

	avg, err := a.repo.Avg(ctx, from, to, categoryID, typeFilter)
	if err != nil {
		return 0, errutils.Wrap(op, err)
	}

	return avg, nil
}

func (a *Analytics) Count(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (int, error) {
	const op = "service.analytics.Count"

	var typeFilter *domain.ItemType
	if itemType != nil {
		t := domain.ItemType(*itemType)
		typeFilter = &t
	}

	count, err := a.repo.Count(ctx, from, to, categoryID, typeFilter)
	if err != nil {
		return 0, errutils.Wrap(op, err)
	}

	return count, nil
}

func (a *Analytics) Median(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error) {
	const op = "service.analytics.Median"

	var typeFilter *domain.ItemType
	if itemType != nil {
		t := domain.ItemType(*itemType)
		typeFilter = &t
	}

	median, err := a.repo.Median(ctx, from, to, categoryID, typeFilter)
	if err != nil {
		return 0, errutils.Wrap(op, err)
	}

	return median, nil
}

func (a *Analytics) PercentileNinetieth(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (float64, error) {
	const op = "service.analytics.PercentileNinetieth"

	var typeFilter *domain.ItemType
	if itemType != nil {
		t := domain.ItemType(*itemType)
		typeFilter = &t
	}

	p90, err := a.repo.PercentileNinetieth(ctx, from, to, categoryID, typeFilter)
	if err != nil {
		return 0, errutils.Wrap(op, err)
	}

	return p90, nil
}
