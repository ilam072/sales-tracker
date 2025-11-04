package service

import (
	"context"
	"errors"
	"github.com/ilam072/sales-tracker/internal/category/repo"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/internal/types/dto"
	"github.com/ilam072/sales-tracker/pkg/errutils"
)

type CategoryRepo interface {
	CreateCategory(ctx context.Context, category domain.Category) (int, error)
	GetCategoryByID(ctx context.Context, id int) (domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, cat domain.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type Category struct {
	repo CategoryRepo
}

func New(repo CategoryRepo) *Category {
	return &Category{repo: repo}
}

func (c *Category) SaveCategory(ctx context.Context, category dto.CreateCategory) (int, error) {
	const op = "service.category.Save"

	domainCategory := domain.Category{
		Name: category.Name,
	}

	ID, err := c.repo.CreateCategory(ctx, domainCategory)
	if err != nil {
		if errors.Is(err, repo.ErrCategoryExists) {
			return 0, errutils.Wrap("failed to create category", domain.ErrCategoryExists)
		}
		return 0, errutils.Wrap("failed to create category", err)
	}

	return ID, nil
}

func (c *Category) GetCategoryByID(ctx context.Context, id int) (dto.GetCategory, error) {
	const op = "service.category.GetByID"

	category, err := c.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrCategoryNotFound) {
			return dto.GetCategory{}, errutils.Wrap(op, domain.ErrCategoryNotFound)
		}
		return dto.GetCategory{}, errutils.Wrap(op, err)
	}

	return dto.GetCategory{ID: category.ID, Name: category.Name}, nil
}

func (c *Category) GetAllCategories(ctx context.Context) (dto.Categories, error) {
	const op = "service.category.GetAll"

	categories, err := c.repo.GetAllCategories(ctx)
	if err != nil {
		return dto.Categories{}, errutils.Wrap(op, err)
	}

	if len(categories) == 0 {
		return dto.Categories{Categories: []dto.GetCategory{}}, nil
	}

	result := make([]dto.GetCategory, 0, len(categories))
	for _, cat := range categories {
		result = append(result, dto.GetCategory{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	return dto.Categories{Categories: result}, nil
}

func (c *Category) UpdateCategory(ctx context.Context, id int, category dto.UpdateCategory) error {
	const op = "service.category.Update"

	domainCategory := domain.Category{ID: id, Name: category.Name}

	if err := c.repo.UpdateCategory(ctx, domainCategory); err != nil {
		if errors.Is(err, repo.ErrCategoryNotFound) {
			return errutils.Wrap(op, domain.ErrCategoryNotFound)
		}
		return errutils.Wrap(op, err)
	}

	return nil
}

func (c *Category) DeleteCategory(ctx context.Context, id int) error {
	const op = "service.category.Delete"

	if err := c.repo.DeleteCategory(ctx, id); err != nil {
		if errors.Is(err, repo.ErrCategoryNotFound) {
			return errutils.Wrap(op, domain.ErrCategoryNotFound)
		}
		return errutils.Wrap(op, err)
	}

	return nil
}
