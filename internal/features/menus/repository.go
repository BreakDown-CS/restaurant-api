package menus

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// menuColumns reads the menu with its category name, joined from categories.
const menuColumns = `m.id, m.name, m.price, m.category_id, c.name AS category, m.is_available, m.created_at, m.updated_at`

const foreignKeyViolation = "23503"

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation
}

func (r *Repository) List(ctx context.Context) ([]Menu, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+menuColumns+`
		FROM menus m
		JOIN categories c ON c.id = m.category_id
		ORDER BY m.created_at DESC, m.id DESC`)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Menu])
}

func (r *Repository) Create(ctx context.Context, in MenuInput) (Menu, error) {
	rows, err := r.db.Query(ctx, `
		WITH written AS (
			INSERT INTO menus (name, price, category_id, is_available)
			VALUES ($1, $2, $3, $4)
			RETURNING *
		)
		SELECT `+menuColumns+`
		FROM written m
		JOIN categories c ON c.id = m.category_id`,
		in.Name, in.Price, in.CategoryID, in.IsAvailable,
	)
	if err != nil {
		return Menu{}, err
	}

	menu, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Menu])
	if isForeignKeyViolation(err) {
		return Menu{}, ErrCategoryNotFound
	}
	return menu, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, in MenuInput) (Menu, error) {
	rows, err := r.db.Query(ctx, `
		WITH written AS (
			UPDATE menus
			SET name = $2, price = $3, category_id = $4, is_available = $5, updated_at = now()
			WHERE id = $1
			RETURNING *
		)
		SELECT `+menuColumns+`
		FROM written m
		JOIN categories c ON c.id = m.category_id`,
		id, in.Name, in.Price, in.CategoryID, in.IsAvailable,
	)
	if err != nil {
		return Menu{}, err
	}

	menu, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Menu])
	if errors.Is(err, pgx.ErrNoRows) {
		return Menu{}, ErrNotFound
	}
	if isForeignKeyViolation(err) {
		return Menu{}, ErrCategoryNotFound
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
