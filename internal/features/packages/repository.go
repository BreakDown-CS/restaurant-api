package packages

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const packageColumns = `id, name, price_per_person, time_limit_minutes, free_below_height_cm,
	overtime_fee, overtime_block_minutes, overtime_grace_minutes,
	waste_fee, waste_unit_grams, is_active, created_at, updated_at`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Package, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+packageColumns+`
		FROM buffet_packages
		ORDER BY price_per_person ASC, name ASC`)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Package])
}
