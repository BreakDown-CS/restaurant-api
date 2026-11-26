package packages

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	decimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestList_ReadsAllColumns(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	_, err = pool.Exec(ctx, `
		INSERT INTO buffet_packages (name, price_per_person, time_limit_minutes,
			free_below_height_cm, overtime_fee, overtime_block_minutes,
			overtime_grace_minutes, waste_fee, waste_unit_grams)
		VALUES ('Test Basic', 299.50, 90, 100, 30, 15, 5, 50, 100)`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM buffet_packages WHERE name = 'Test Basic'`)
	})

	got, err := NewRepository(pool).List(ctx)

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "Test Basic", got[0].Name)
	require.Equal(t, "299.5", got[0].PricePerPerson.String())
	require.Equal(t, 90, got[0].TimeLimitMinutes)
	require.Equal(t, 5, got[0].OvertimeGraceMinutes)
	require.True(t, got[0].IsActive)
}

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	pool, err := pgxpool.New(context.Background(), url)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	return pool
}

func sampleInput(name string) PackageInput {
	return PackageInput{
		Name:                 name,
		PricePerPerson:       decimal.RequireFromString("399.00"),
		TimeLimitMinutes:     90,
		FreeBelowHeightCm:    100,
		OvertimeFee:          decimal.RequireFromString("30"),
		OvertimeBlockMinutes: 15,
		OvertimeGraceMinutes: 5,
		WasteFee:             decimal.RequireFromString("50"),
		WasteUnitGrams:       100,
		IsActive:             true,
	}
}

func TestCreateGetUpdateDelete(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	repo := NewRepository(pool)

	created, err := repo.Create(ctx, sampleInput("Test Premium"))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM buffet_packages WHERE id = $1`, created.ID)
	})
	require.True(t, created.PricePerPerson.Equal(decimal.RequireFromString("399")))

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "Test Premium", got.Name)

	in := sampleInput("Test Premium Plus")
	in.TimeLimitMinutes = 120
	updated, err := repo.Update(ctx, created.ID, in)
	require.NoError(t, err)
	require.Equal(t, "Test Premium Plus", updated.Name)
	require.Equal(t, 120, updated.TimeLimitMinutes)

	require.NoError(t, repo.Delete(ctx, created.ID))
	_, err = repo.GetByID(ctx, created.ID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestCreate_DuplicateName(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	repo := NewRepository(pool)

	first, err := repo.Create(ctx, sampleInput("Test Duplicate"))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM buffet_packages WHERE id = $1`, first.ID)
	})

	_, err = repo.Create(ctx, sampleInput("Test Duplicate"))
	require.ErrorIs(t, err, ErrNameTaken)
}

func TestNotFound(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	repo := NewRepository(pool)
	missing := uuid.New()

	_, err := repo.GetByID(ctx, missing)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = repo.Update(ctx, missing, sampleInput("Test Missing"))
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, repo.Delete(ctx, missing), ErrNotFound)
}
