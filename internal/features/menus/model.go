// Package menus manages the restaurant's menu items.
package menus

import (
	"errors"
	"time"

	"github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

var (
	ErrNotFound         = errors.New("menu not found")
	ErrCategoryNotFound = errors.New("category not found")
)

type Menu struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Price           decimal.Decimal `json:"price"`
	CategoryID      uuid.UUID       `json:"categoryId"`
	Category        string          `json:"category"`
	IsAvailable     bool            `json:"isAvailable"`
	MinPackageLevel *int            `json:"minPackageLevel"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type MenuInput struct {
	Name            string           `json:"name" binding:"required,max=100"`
	Price           *decimal.Decimal `json:"price" binding:"required"`
	CategoryID      uuid.UUID        `json:"categoryId" binding:"required"`
	IsAvailable     bool             `json:"isAvailable"`
	MinPackageLevel *int             `json:"minPackageLevel"`
}
