package menus

import "context"

type Usecase struct {
	repo *Repository
}

func NewUsecase(repo *Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context) ([]Menu, error) {
	return u.repo.List(ctx)
}
