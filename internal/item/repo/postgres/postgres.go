package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ilam072/sales-tracker/internal/item/repo"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/pkg/errutils"
	"github.com/wb-go/wbf/dbpg"
	"strings"
	"time"
)

type ItemRepo struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

func (r *ItemRepo) CreateItem(ctx context.Context, item domain.Item) (int, error) {
	query := `
        INSERT INTO items (category_id, type, amount, description, transaction_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
    `
	var id int
	if err := r.db.QueryRowContext(ctx, query,
		item.CategoryId,
		item.Type,
		item.Amount,
		item.Description,
		item.TransactionDate,
	).Scan(&id); err != nil {
		return 0, errutils.Wrap("failed to create item", err)
	}
	return id, nil
}

func (r *ItemRepo) GetItemByID(ctx context.Context, id int) (domain.Item, error) {
	query := `
        SELECT id, category_id, type, amount, description, created_at, transaction_date
        FROM items
        WHERE id = $1;
    `
	var item domain.Item
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.Id,
		&item.CategoryId,
		&item.Type,
		&item.Amount,
		&item.Description,
		&item.CreatedAt,
		&item.TransactionDate,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Item{}, errutils.Wrap("failed to get item by id", repo.ErrItemNotFound)
		}
		return domain.Item{}, errutils.Wrap("failed to get item by id", err)
	}
	return item, nil
}

func (r *ItemRepo) GetAllItems(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) ([]domain.Item, error) {
	query := `
        SELECT id, category_id, type, amount, description, created_at, transaction_date
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

	query += " ORDER BY transaction_date DESC;"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errutils.Wrap("failed to get all items", err)
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(
			&item.Id,
			&item.CategoryId,
			&item.Type,
			&item.Amount,
			&item.Description,
			&item.CreatedAt,
			&item.TransactionDate,
		); err != nil {
			return nil, errutils.Wrap("failed to scan item", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ItemRepo) UpdateItem(ctx context.Context, item domain.Item) error {
	query := `
        UPDATE items
        SET category_id = $1,
            type = $2,
            amount = $3,
            description = $4,
            transaction_date = $5
        WHERE id = $6;
    `
	res, err := r.db.ExecContext(ctx, query,
		item.CategoryId,
		item.Type,
		item.Amount,
		item.Description,
		item.TransactionDate,
		item.Id,
	)
	if err != nil {
		return errutils.Wrap("failed to update item", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errutils.Wrap("failed to get affected rows number", err)
	}

	if rows == 0 {
		return repo.ErrItemNotFound
	}

	return nil
}

func (r *ItemRepo) DeleteItem(ctx context.Context, id int) error {
	query := `DELETE FROM items WHERE id = $1;`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return errutils.Wrap("failed to delete item", err)
	}
	return nil
}
