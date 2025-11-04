package domain

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrItemNotFound     = errors.New("item not found")
)
