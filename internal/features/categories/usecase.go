package categories

import "context"

type repository interface {
	List(ctx context.Context) ([]Category, error)
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
