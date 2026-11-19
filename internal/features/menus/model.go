// Package menus manages the restaurant's menu items.
package menus

import (
	"time"

	"github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

type Menu struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Category    string          `json:"category"`
	IsAvailable bool            `json:"isAvailable"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
