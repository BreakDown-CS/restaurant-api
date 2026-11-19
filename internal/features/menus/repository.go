package menus

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Menu, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, price, category, is_available, created_at, updated_at
		FROM menus
		ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Menu])
}
