package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ilam072/sales-tracker/internal/category/repo"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/pkg/errutils"
	"github.com/lib/pq"
	"github.com/wb-go/wbf/dbpg"
)

type CategoryRepo struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) CreateCategory(ctx context.Context, category domain.Category) (int, error) {
	query := `
        INSERT INTO categories (name)
        VALUES ($1)
        RETURNING id;
    `

	var id int
	if err := r.db.QueryRowContext(ctx, query, category.Name).Scan(&id); err != nil {
		if isUniqueViolation(err) {
			return 0, errutils.Wrap("failed to create category", repo.ErrCategoryExists)
		}
		return 0, errutils.Wrap("failed to create category", err)
	}

	return id, nil
}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, id int) (domain.Category, error) {
	query := `
        SELECT id, name, created_at
        FROM categories
        WHERE id = $1;
    `

	var category domain.Category
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, errutils.Wrap("failed to get category by id", repo.ErrCategoryNotFound)
		}
		return domain.Category{}, errutils.Wrap("failed to get category by id", err)
	}

	return category, nil
}

func (r *CategoryRepo) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	query := `
        SELECT id, name, created_at
        FROM categories
        ORDER BY created_at DESC;
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errutils.Wrap("failed to get all categories", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.CreatedAt,
		); err != nil {
			return nil, errutils.Wrap("failed to scan category", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, cat domain.Category) error {
	query := `
        UPDATE categories
        SET name = $1
        WHERE id = $2;
    `

	res, err := r.db.ExecContext(ctx, query, cat.Name, cat.ID)
	if err != nil {
		return errutils.Wrap("failed to update category", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errutils.Wrap("failed to get affected rows number", err)
	}

	if rows == 0 {
		return repo.ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1;`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errutils.Wrap("failed to delete category", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errutils.Wrap("failed to get affected rows number", err)
	}

	if rows == 0 {
		return repo.ErrCategoryNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
