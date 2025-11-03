package dto

import "time"

type CreateItem struct {
	CategoryId      int       `json:"category_id"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date,omitempty"`
}

type GetItem struct {
	CategoryId      int       `json:"category_id"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}

type UpdateItem struct {
	CategoryId      int       `json:"category_id"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date,omitempty"`
}

type Items struct {
	Items []GetItem `json:"items"`
}
