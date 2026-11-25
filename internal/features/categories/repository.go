package categories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const categoryColumns = `id, name, created_at, updated_at`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+categoryColumns+`
		FROM categories
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Category])
}

const (
	uniqueViolation     = "23505"
	foreignKeyViolation = "23503"
)

func hasPgCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}

func (r *Repository) Create(ctx context.Context, in CategoryInput) (Category, error) {
	rows, err := r.db.Query(ctx, `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING `+categoryColumns,
		in.Name,
	)
	if err != nil {
		return Category{}, err
	}

	cat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Category])
	if hasPgCode(err, uniqueViolation) {
		return Category{}, ErrNameTaken
	}

	return cat, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, in CategoryInput) (Category, error) {
	rows, err := r.db.Query(ctx, `
		UPDATE categories
		SET name = $2, updated_at = now()
		WHERE id = $1
		RETURNING `+categoryColumns,
		id, in.Name,
	)

	if err != nil {
		return Category{}, err
	}

	cat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Category])
	if errors.Is(err, pgx.ErrNoRows) {
		return Category{}, ErrNotFound
	}

	if hasPgCode(err, uniqueViolation) {
		return Category{}, ErrNameTaken
	}

	return cat, err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)

	if hasPgCode(err, foreignKeyViolation) {
		return ErrInUse
	}

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
