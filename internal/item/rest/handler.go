package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilam072/sales-tracker/internal/response"
	"github.com/ilam072/sales-tracker/internal/types/domain"
	"github.com/ilam072/sales-tracker/internal/types/dto"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
	"strconv"
	"time"
)

type Item interface {
	CreateItem(ctx context.Context, item dto.CreateItem) (int, error)
	GetItemByID(ctx context.Context, id int) (dto.GetItem, error)
	GetAllItems(ctx context.Context, from, to *time.Time, categoryID *int, itemType *string) (dto.Items, error)
	UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error
	DeleteItem(ctx context.Context, id int) error
}

type Validator interface {
	Validate(i interface{}) error
}

type ItemHandler struct {
	item      Item
	validator Validator
}

func NewItemHandler(item Item, validator Validator) *ItemHandler {
	return &ItemHandler{item: item, validator: validator}
}

func (h *ItemHandler) CreateItem(c *ginext.Context) {
	var item dto.CreateItem
	if err := c.BindJSON(&item); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind item JSON")
		response.Error("invalid request body").WriteJSON(c, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(item); err != nil {
		zlog.Logger.Error().Err(err).Msg("validation error")
		response.Error(fmt.Sprintf("validation error: %s", err.Error())).WriteJSON(c, http.StatusBadRequest)
		return
	}

	// Если дата транзакции не указана — ставим текущую
	if item.TransactionDate.IsZero() {
		item.TransactionDate = time.Now().UTC()
	}

	ID, err := h.item.CreateItem(c.Request.Context(), item)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("failed to create item: category not found")
			response.Error("category not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to create item")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusCreated, ginext.H{"item_id": ID})
}

func (h *ItemHandler) GetItemByID(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid item id param")
		response.Error("invalid item id").WriteJSON(c, http.StatusBadRequest)
		return
	}

	item, err := h.item.GetItemByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			zlog.Logger.Error().Err(err).Msg("item not found")
			response.Error("item not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to get item by id")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"item": item})
}

func (h *ItemHandler) GetAllItems(c *ginext.Context) {
	var (
		fromStr       = c.Query("from")
		toStr         = c.Query("to")
		categoryIDStr = c.Query("category_id")
		typeStr       = c.Query("type")
	)

	var (
		from, to   *time.Time
		categoryID *int
		itemType   *string
	)

	if fromStr != "" {
		t, err := time.Parse(time.DateOnly, fromStr)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("invalid from date format")
			response.Error("invalid 'from' format, expected YYYY-MM-DD").WriteJSON(c, http.StatusBadRequest)
			return
		}
		from = &t
	}

	if toStr != "" {
		t, err := time.Parse(time.DateOnly, toStr)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("invalid to date format")
			response.Error("invalid 'to' format, expected YYYY-MM-DD").WriteJSON(c, http.StatusBadRequest)
			return
		}
		to = &t
	}

	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("invalid category_id param")
			response.Error("invalid 'category_id', must be an integer").WriteJSON(c, http.StatusBadRequest)
			return
		}
		categoryID = &id
	}

	if typeStr != "" {
		itemType = &typeStr
	}

	items, err := h.item.GetAllItems(c.Request.Context(), from, to, categoryID, itemType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get all items")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, items)
}

func (h *ItemHandler) UpdateItem(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid item id param")
		response.Error("invalid item id, must be an integer").WriteJSON(c, http.StatusBadRequest)
		return
	}

	var item dto.UpdateItem
	if err := c.BindJSON(&item); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind item JSON")
		response.Error("invalid request body").WriteJSON(c, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(item); err != nil {
		zlog.Logger.Error().Err(err).Msg("validation error")
		response.Error(fmt.Sprintf("validation error: %s", err.Error())).WriteJSON(c, http.StatusBadRequest)
		return
	}

	if item.TransactionDate.IsZero() {
		item.TransactionDate = time.Now().UTC()
	}

	err = h.item.UpdateItem(c.Request.Context(), id, item)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			response.Error("item not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to update item")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"message": "item successfully updated"})
}

func (h *ItemHandler) DeleteItem(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid item id param")
		response.Error("invalid item id, must be an integer").WriteJSON(c, http.StatusBadRequest)
		return
	}

	err = h.item.DeleteItem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			response.Error("item not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to delete item")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"message": "item successfully deleted"})
}
