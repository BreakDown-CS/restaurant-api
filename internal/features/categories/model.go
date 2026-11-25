// Package categories manages the menu categories shown in the admin dropdown
package categories

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var (
	ErrNotFound  = errors.New("category not found")
	ErrNameTaken = errors.New("category name already taken")
	ErrInUse     = errors.New("category still has menus")
)

type CategoryInput struct {
	Name string `json:"name"`
}
