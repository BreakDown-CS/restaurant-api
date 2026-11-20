package menus

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

var maxPrice = decimal.RequireFromString("99999999.99")

type Usecase struct {
	repo *Repository
}

func NewUsecase(repo *Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context) ([]Menu, error) {
	return u.repo.List(ctx)
}

func (u *Usecase) Create(ctx context.Context, in MenuInput) (Menu, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Category = strings.TrimSpace(in.Category)

	if err := validateInput(in); err != nil {
		return Menu{}, err
	}

	return u.repo.Create(ctx, in)
}

func validateInput(in MenuInput) error {
	fields := map[string]string{}

	if in.Name == "" {
		fields["name"] = "required"
	}

	if in.Category == "" {
		fields["category"] = "required"
	}

	switch {
	case in.Price == nil:
		fields["price"] = "required"
	case in.Price.IsNegative():
		fields["price"] = "must be 0 or more"
	case !in.Price.Equal(in.Price.Round(2)):
		fields["price"] = "must have at most 2 decimal places"
	case in.Price.GreaterThan(maxPrice):
		fields["price"] = "must be at most 99999999.99"
	}

	if len(fields) > 0 {
		return &validation.Error{Fields: fields}
	}

	return nil
}
