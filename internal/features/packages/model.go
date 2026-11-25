// Package packages manages buffet packages and the menus each one includes.
package packages

import (
	"time"

	"github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

type Package struct {
	ID                   uuid.UUID       `json:"id"`
	Name                 string          `json:"name"`
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
