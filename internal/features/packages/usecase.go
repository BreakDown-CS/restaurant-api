package packages

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

var maxMoney = decimal.RequireFromString("99999999.99")

type repository interface {
	List(ctx context.Context) ([]Package, error)
	GetByID(ctx context.Context, id uuid.UUID) (Package, error)
	Create(ctx context.Context, in PackageInput) (Package, error)
	Update(ctx context.Context, id uuid.UUID, in PackageInput) (Package, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	repo repository
}

func NewUsecase(repo repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context) ([]Package, error) {
	return u.repo.List(ctx)
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (Package, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Usecase) Create(ctx context.Context, in PackageInput) (Package, error) {
	in = cleanInput(in)

	if err := validateInput(in); err != nil {
		return Package{}, err
	}

	return u.repo.Create(ctx, in)
}

func (u *Usecase) Update(ctx context.Context, id uuid.UUID, in PackageInput) (Package, error) {
	in = cleanInput(in)
	if err := validateInput(in); err != nil {
		return Package{}, err
	}

	return u.repo.Update(ctx, id, in)
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func cleanInput(in PackageInput) PackageInput {
	in.Name = strings.TrimSpace(in.Name)
	if in.IsActive == nil {
		in.IsActive = new(true)
	}
	return in
}

func validateInput(in PackageInput) error {
	fields := map[string]string{}

	if in.Name == "" {
		fields["name"] = "required"
	}

	if in.PricePerPerson.Sign() <= 0 {
		fields["pricePerPerson"] = "must be more than 0"
	} else if msg := checkMoney(in.PricePerPerson); msg != "" {
		fields["pricePerPerson"] = msg
	}

	if msg := checkMoney(in.OvertimeFee); msg != "" {
		fields["overtimeFee"] = msg
	}

	if msg := checkMoney(in.WasteFee); msg != "" {
		fields["wasteFee"] = msg
	}

	mustBePositive(fields, "level", in.Level)
	mustBePositive(fields, "timeLimitMinutes", in.TimeLimitMinutes)
	mustBePositive(fields, "overtimeBlockMinutes", in.OvertimeBlockMinutes)
	mustBePositive(fields, "wasteUnitGrams", in.WasteUnitGrams)

	mustNotBeNegative(fields, "freeBelowHeightCm", in.FreeBelowHeightCm)
	mustNotBeNegative(fields, "overtimeGraceMinutes", in.OvertimeGraceMinutes)

	if len(fields) > 0 {
		return &validation.Error{Fields: fields}
	}

	return nil
}

func checkMoney(d decimal.Decimal) string {
	switch {
	case d.IsNegative():
		return "must be 0 or more"
	case !d.Equal(d.Round(2)):
		return "must have at most 2 decimal places"
	case d.GreaterThan(maxMoney):
		return "must be at most 99999999.99"
	}
	return ""
}

func mustBePositive(fields map[string]string, name string, v int) {
	if v <= 0 {
		fields[name] = "must be more than 0"
	}
}

func mustNotBeNegative(fields map[string]string, name string, v int) {
	if v < 0 {
		fields[name] = "must be 0 or more"
	}
}
