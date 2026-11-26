package packages

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	List(ctx context.Context) ([]Package, error)
	GetByID(ctx context.Context, id uuid.UUID) (Package, error)
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
