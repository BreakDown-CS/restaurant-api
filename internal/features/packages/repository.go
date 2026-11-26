package packages

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const packageColumns = `id, name, level, price_per_person, time_limit_minutes, free_below_height_cm,
	overtime_fee, overtime_block_minutes, overtime_grace_minutes,
	waste_fee, waste_unit_grams, is_active, created_at, updated_at`

const uniqueViolation = "23505"

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolation
}

func (r *Repository) List(ctx context.Context) ([]Package, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+packageColumns+`
		FROM buffet_packages
		ORDER BY level ASC, price_per_person ASC, name ASC`)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Package])
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (Package, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+packageColumns+`
		FROM buffet_packages
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return Package{}, err
	}

	pkg, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Package])
	if errors.Is(err, pgx.ErrNoRows) {
		return Package{}, ErrNotFound
	}

	return pkg, err
}

func (r *Repository) Create(ctx context.Context, in PackageInput) (Package, error) {
	rows, err := r.db.Query(ctx, `
		INSERT INTO buffet_packages (name, price_per_person, time_limit_minutes,
			free_below_height_cm, overtime_fee, overtime_block_minutes,
			overtime_grace_minutes, waste_fee, waste_unit_grams, is_active, level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING `+packageColumns,
		in.Name, in.PricePerPerson, in.TimeLimitMinutes, in.FreeBelowHeightCm, in.OvertimeFee,
		in.OvertimeBlockMinutes, in.OvertimeGraceMinutes, in.WasteFee, in.WasteUnitGrams, in.IsActive,
		in.Level,
	)

	if err != nil {
		return Package{}, err
	}

	pkg, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Package])
	if isUniqueViolation(err) {
		return Package{}, ErrNameTaken
	}

	return pkg, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, in PackageInput) (Package, error) {
	rows, err := r.db.Query(ctx, `
		UPDATE buffet_packages
		SET name = $2, price_per_person = $3, time_limit_minutes = $4, free_below_height_cm = $5,
			overtime_fee = $6, overtime_block_minutes = $7, overtime_grace_minutes = $8, waste_fee = $9,
			waste_unit_grams = $10, is_active = $11, level = $12, updated_at = now()
		WHERE id = $1
		RETURNING `+packageColumns,
		id, in.Name, in.PricePerPerson, in.TimeLimitMinutes, in.FreeBelowHeightCm, in.OvertimeFee,
		in.OvertimeBlockMinutes, in.OvertimeGraceMinutes, in.WasteFee, in.WasteUnitGrams, in.IsActive,
		in.Level,
	)

	if err != nil {
		return Package{}, err
	}

	pkg, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Package])
	if errors.Is(err, pgx.ErrNoRows) {
		return Package{}, ErrNotFound
	}

	if isUniqueViolation(err) {
		return Package{}, ErrNameTaken
	}

	return pkg, err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM buffet_packages WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
