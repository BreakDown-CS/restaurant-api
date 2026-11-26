// Package packages manages buffet packages and the menus each one includes.
package packages

import (
	"errors"
	"time"

	"github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

type Package struct {
	ID                   uuid.UUID       `json:"id"`
	Name                 string          `json:"name"`
	Level                int             `json:"level"`
	PricePerPerson       decimal.Decimal `json:"pricePerPerson"`
	TimeLimitMinutes     int             `json:"timeLimitMinutes"`
	FreeBelowHeightCm    int             `json:"freeBelowHeightCm"`
	OvertimeFee          decimal.Decimal `json:"overtimeFee"`
	OvertimeBlockMinutes int             `json:"overtimeBlockMinutes"`
	OvertimeGraceMinutes int             `json:"overtimeGraceMinutes"`
	WasteFee             decimal.Decimal `json:"wasteFee"`
	WasteUnitGrams       int             `json:"wasteUnitGrams"`
	IsActive             bool            `json:"isActive"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
}

var (
	ErrNotFound  = errors.New("buffet package not found")
	ErrNameTaken = errors.New("buffet package name already taken")
)

type PackageInput struct {
	Name                 string          `json:"name" binding:"required,max=100"`
	Level                int             `json:"level"`
	PricePerPerson       decimal.Decimal `json:"pricePerPerson"`
	TimeLimitMinutes     int             `json:"timeLimitMinutes"`
	FreeBelowHeightCm    int             `json:"freeBelowHeightCm"`
	OvertimeFee          decimal.Decimal `json:"overtimeFee"`
	OvertimeBlockMinutes int             `json:"overtimeBlockMinutes"`
	OvertimeGraceMinutes int             `json:"overtimeGraceMinutes"`
	WasteFee             decimal.Decimal `json:"wasteFee"`
	WasteUnitGrams       int             `json:"wasteUnitGrams"`
	IsActive             *bool           `json:"isActive"`
}
