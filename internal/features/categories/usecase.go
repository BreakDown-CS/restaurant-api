package categories

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

const maxNameLength = 50

type repository interface {
	List(ctx context.Context) ([]Category, error)
	Create(ctx context.Context, in CategoryInput) (Category, error)
	Update(ctx context.Context, id uuid.UUID, in CategoryInput) (Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	repo repository
}

func NewUsecase(repo repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context) ([]Category, error) {
	return u.repo.List(ctx)
}

func (u *Usecase) Create(ctx context.Context, in CategoryInput) (Category, error) {
	in = cleanInput(in)
	if err := validateInput(in); err != nil {
		return Category{}, err
	}

	return u.repo.Create(ctx, in)
}

func (u *Usecase) Update(ctx context.Context, id uuid.UUID, in CategoryInput) (Category, error) {
	in = cleanInput(in)
	if err := validateInput(in); err != nil {
		return Category{}, err
	}

	return u.repo.Update(ctx, id, in)
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func cleanInput(in CategoryInput) CategoryInput {
	in.Name = strings.TrimSpace(in.Name)
	return in
}

func validateInput(in CategoryInput) error {
	fields := map[string]string{}

	switch {
	case in.Name == "":
		fields["name"] = "required"
	case utf8.RuneCountInString(in.Name) > maxNameLength:
		fields["name"] = "must be at most 50 characters"
	}

	if len(fields) > 0 {
		return &validation.Error{Fields: fields}
	}

	return nil
}
