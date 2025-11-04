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
)

type Category interface {
	SaveCategory(ctx context.Context, category dto.CreateCategory) (int, error)
	GetCategoryByID(ctx context.Context, id int) (dto.GetCategory, error)
	GetAllCategories(ctx context.Context) (dto.Categories, error)
	UpdateCategory(ctx context.Context, id int, category dto.UpdateCategory) error
	DeleteCategory(ctx context.Context, id int) error
}

type Validator interface {
	Validate(i interface{}) error
}

type CategoryHandler struct {
	category  Category
	validator Validator
}

func NewCategoryHandler(category Category, validator Validator) *CategoryHandler {
	return &CategoryHandler{category: category, validator: validator}
}

func (h *CategoryHandler) CreateCategory(c *ginext.Context) {
	var category dto.CreateCategory
	if err := c.BindJSON(&category); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind category JSON")
		response.Error("invalid request body").WriteJSON(c, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(category); err != nil {
		zlog.Logger.Error().Err(err).Msg("validation error")
		response.Error(fmt.Sprintf("validation error: %s", err.Error())).WriteJSON(c, http.StatusBadRequest)
		return
	}

	ID, err := h.category.SaveCategory(c.Request.Context(), category)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryExists) {
			zlog.Logger.Error().Err(err).Msg("failed to create category")
			response.Error("category with this name already exists").WriteJSON(c, http.StatusConflict)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to create category")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusCreated, ginext.H{"category_id": ID})
}

func (h *CategoryHandler) GetCategoryByID(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid category id param")
		response.Error("invalid category id").WriteJSON(c, http.StatusBadRequest)
		return
	}

	category, err := h.category.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("category not found")
			response.Error("category not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to get category by id")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, ginext.H{"category": category})
}

func (h *CategoryHandler) GetAllCategories(c *ginext.Context) {
	categories, err := h.category.GetAllCategories(c.Request.Context())
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get all categories")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Raw(c, http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid category id param")
		response.Error("invalid category id").WriteJSON(c, http.StatusBadRequest)
		return
	}

	var category dto.UpdateCategory
	if err := c.BindJSON(&category); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind category JSON")
		response.Error("invalid request body").WriteJSON(c, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(category); err != nil {
		zlog.Logger.Error().Err(err).Msg("validation error")
		response.Error(fmt.Sprintf("validation error: %s", err.Error())).WriteJSON(c, http.StatusBadRequest)
		return
	}

	if err := h.category.UpdateCategory(c.Request.Context(), id, category); err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("category not found")
			response.Error("category not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to update category")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Success("category updated successfully").WriteJSON(c, http.StatusOK)
}

func (h *CategoryHandler) DeleteCategory(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid category id param")
		response.Error("invalid category id").WriteJSON(c, http.StatusBadRequest)
		return
	}

	if err := h.category.DeleteCategory(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("category not found")
			response.Error("category not found").WriteJSON(c, http.StatusNotFound)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to delete category")
		response.Error("internal server error, try again later").WriteJSON(c, http.StatusInternalServerError)
		return
	}

	response.Success("category deleted successfully").WriteJSON(c, http.StatusOK)
}
