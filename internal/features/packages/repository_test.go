package packages

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
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
