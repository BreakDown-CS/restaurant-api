package menus

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

var maxPrice = decimal.RequireFromString("99999999.99")

type repository interface {
	List(ctx context.Context) ([]Menu, error)
	Create(ctx context.Context, in MenuInput) (Menu, error)
	Update(ctx context.Context, id uuid.UUID, in MenuInput) (Menu, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	repo repository
}

func NewUsecase(repo repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context) ([]Menu, error) {
	return u.repo.List(ctx)
}

func (u *Usecase) Create(ctx context.Context, in MenuInput) (Menu, error) {
	in = cleanInput(in)

	if err := validateInput(in); err != nil {
		return Menu{}, err
	}

	return u.repo.Create(ctx, in)
}

func (u *Usecase) Update(ctx context.Context, id uuid.UUID, in MenuInput) (Menu, error) {
	in = cleanInput(in)
	if err := validateInput(in); err != nil {
		return Menu{}, err
	}

	return u.repo.Update(ctx, id, in)
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func cleanInput(in MenuInput) MenuInput {
	in.Name = strings.TrimSpace(in.Name)
	return in
}

func validateInput(in MenuInput) error {
	fields := map[string]string{}

	if in.Name == "" {
		fields["name"] = "required"
	}

	if in.CategoryID == uuid.Nil {
		fields["categoryId"] = "required"
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
