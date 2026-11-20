package menus

import (
	"context"
	"errors"

	"github.com/google/uuid"
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

func (r *Repository) Update(ctx context.Context, id uuid.UUID, in MenuInput) (Menu, error) {
	rows, err := r.db.Query(ctx, `
		UPDATE menus
		SET name = $2, price = $3, category = $4, is_available = $5, updated_at = now()
		WHERE id = $1
		RETURNING `+menuColumns,
		id, in.Name, in.Price, in.Category, in.IsAvailable,
	)
	if err != nil {
		return Menu{}, err
	}

	menu, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Menu])
	if errors.Is(err, pgx.ErrNoRows) {

		return Menu{}, ErrNotFound
	}
	return menu, err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM menus WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
