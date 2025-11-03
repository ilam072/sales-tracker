package domain

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrItemNotFound     = errors.New("item not found")
)
