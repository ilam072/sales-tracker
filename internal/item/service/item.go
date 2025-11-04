package service

import (
	"context"
	"errors"
	"github.com/ilam072/sales-tracker/internal/item/repo"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/internal/types/dto"
	"github.com/ilam072/sales-tracker/pkg/errutils"
	"time"
)

type ItemRepo interface {
	CreateItem(ctx context.Context, item domain.Item) (int, error)
	GetItemByID(ctx context.Context, id int) (domain.Item, error)
	GetAllItems(ctx context.Context, from, to *time.Time, categoryID *int, itemType *domain.ItemType) ([]domain.Item, error)
	UpdateItem(ctx context.Context, item domain.Item) error
	DeleteItem(ctx context.Context, id int) error
}

type Item struct {
	repo ItemRepo
}

func New(repo ItemRepo) *Item {
	return &Item{repo: repo}
}

func (i *Item) CreateItem(ctx context.Context, item dto.CreateItem) (int, error) {
	const op = "service.item.Create"

	transactionDate := item.TransactionDate
	if transactionDate.IsZero() {
		transactionDate = time.Now()
	}

	domainItem := domain.Item{
		CategoryId:      item.CategoryId,
		Type:            domain.ItemType(item.Type),
		Amount:          item.Amount,
		Description:     item.Description,
		TransactionDate: transactionDate,
	}

	id, err := i.repo.CreateItem(ctx, domainItem)
	if err != nil {
		return 0, errutils.Wrap(op, err)
	}

	return id, nil
}

func (i *Item) GetItemByID(ctx context.Context, id int) (dto.GetItem, error) {
	const op = "service.item.GetByID"

	item, err := i.repo.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrItemNotFound) {
			return dto.GetItem{}, errutils.Wrap(op, domain.ErrItemNotFound)
		}
		return dto.GetItem{}, errutils.Wrap(op, err)
	}

	return dto.GetItem{
		CategoryId:      item.CategoryId,
		Type:            string(item.Type),
		Amount:          item.Amount,
		Description:     item.Description,
		TransactionDate: item.TransactionDate,
	}, nil
}

func (i *Item) GetAllItems(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (dto.Items, error) {
	const op = "service.item.GetAll"

	var typeFilter *domain.ItemType
	if itemType != nil {
		t := domain.ItemType(*itemType)
		typeFilter = &t
	}

	items, err := i.repo.GetAllItems(ctx, from, to, categoryID, typeFilter)
	if err != nil {
		return dto.Items{}, errutils.Wrap(op, err)
	}

	if len(items) == 0 {
		return dto.Items{Items: []dto.GetItem{}}, err
	}

	result := make([]dto.GetItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.GetItem{
			CategoryId:      item.CategoryId,
			Type:            string(item.Type),
			Amount:          item.Amount,
			Description:     item.Description,
			TransactionDate: item.TransactionDate,
		})
	}

	return dto.Items{Items: result}, nil
}

func (i *Item) UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error {
	const op = "service.item.Update"

	transactionDate := item.TransactionDate
	if transactionDate.IsZero() {
		transactionDate = time.Now()
	}

	domainItem := domain.Item{
		Id:              id,
		CategoryId:      item.CategoryId,
		Type:            domain.ItemType(item.Type),
		Amount:          item.Amount,
		Description:     item.Description,
		TransactionDate: transactionDate,
	}

	if err := i.repo.UpdateItem(ctx, domainItem); err != nil {
		if errors.Is(err, repo.ErrItemNotFound) {
			return errutils.Wrap(op, domain.ErrItemNotFound)
		}
		return errutils.Wrap(op, err)
	}

	return nil
}

func (i *Item) DeleteItem(ctx context.Context, id int) error {
	const op = "service.item.Delete"

	if err := i.repo.DeleteItem(ctx, id); err != nil {
		if errors.Is(err, repo.ErrItemNotFound) {
			return errutils.Wrap(op, domain.ErrItemNotFound)
		}
		return errutils.Wrap(op, err)
	}

	return nil
}
