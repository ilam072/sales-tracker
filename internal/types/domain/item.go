package domain

import (
	"time"
)

type ItemType string

type Item struct {
	Id              int
	CategoryId      int
	Type            ItemType
	Amount          float64
	Description     string
	CreatedAt       time.Time
	TransactionDate time.Time
}
