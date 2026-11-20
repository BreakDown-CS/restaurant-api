package menus

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const menuColumns = `id, name, price, category, is_available, created_at, updated_at`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Menu, error) {

	rows, err := r.db.Query(ctx, `
		SELECT `+menuColumns+`
		FROM menus
		ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Menu])
}

func (r *Repository) Create(ctx context.Context, in MenuInput) (Menu, error) {
	rows, err := r.db.Query(ctx, `
		INSERT INTO menus (name, price, category, is_available)
		VALUES ($1, $2, $3, $4)
		RETURNING `+menuColumns,
		in.Name, in.Price, in.Category, in.IsAvailable,
	)
	if err != nil {
		return Menu{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Menu])
}
