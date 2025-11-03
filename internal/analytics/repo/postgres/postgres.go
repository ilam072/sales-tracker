package postgres

import (
	"context"
	"fmt"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/pkg/errutils"
	"github.com/wb-go/wbf/dbpg"
	"strings"
	"time"
)

type AnalyticsRepo struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

func (a *AnalyticsRepo) Sum(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error) {
	query := `
        SELECT COALESCE(SUM(amount), 0)
        FROM items
    `

	var conditions []string
	var args []any

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", len(args)+1))
		args = append(args, *from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", len(args)+1))
		args = append(args, *to)
	}

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *categoryID)
	}

	if itemType != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, *itemType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var sum float64
	if err := a.db.QueryRowContext(ctx, query, args...).Scan(&sum); err != nil {
		return 0, errutils.Wrap("failed to calculate sum", err)
	}

	return sum, nil
}

func (a *AnalyticsRepo) Avg(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error) {
	query := `
        SELECT COALESCE(AVG(amount), 0)
        FROM items
    `

	var conditions []string
	var args []any

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", len(args)+1))
		args = append(args, *from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", len(args)+1))
		args = append(args, *to)
	}

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *categoryID)
	}

	if itemType != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, *itemType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var avg float64
	if err := a.db.QueryRowContext(ctx, query, args...).Scan(&avg); err != nil {
		return 0, errutils.Wrap("failed to calculate average", err)
	}

	return avg, nil
}

func (a *AnalyticsRepo) Count(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM items
    `

	var conditions []string
	var args []any

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", len(args)+1))
		args = append(args, *from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", len(args)+1))
		args = append(args, *to)
	}

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *categoryID)
	}

	if itemType != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, *itemType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	if err := a.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, errutils.Wrap("failed to count items", err)
	}

	return count, nil
}

func (a *AnalyticsRepo) Median(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error) {
	query := `
        SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount)
        FROM items
    `

	var conditions []string
	var args []any

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", len(args)+1))
		args = append(args, *from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", len(args)+1))
		args = append(args, *to)
	}

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *categoryID)
	}

	if itemType != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, *itemType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var median float64
	if err := a.db.QueryRowContext(ctx, query, args...).Scan(&median); err != nil {
		return 0, errutils.Wrap("failed to calculate median", err)
	}

	return median, nil
}

func (a *AnalyticsRepo) PercentileNinetieth(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) (float64, error) {
	query := `
        SELECT PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount)
        FROM items
    `

	var conditions []string
	var args []any

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", len(args)+1))
		args = append(args, *from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", len(args)+1))
		args = append(args, *to)
	}

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *categoryID)
	}

	if itemType != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, *itemType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var p90 float64
	if err := a.db.QueryRowContext(ctx, query, args...).Scan(&p90); err != nil {
		return 0, errutils.Wrap("failed to calculate 90th percentile", err)
	}

	return p90, nil
}
