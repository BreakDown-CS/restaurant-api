package categories

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestList_ReadsSeededCategories(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	pool, err := pgxpool.New(context.Background(), url)
	require.NoError(t, err)
	defer pool.Close()

	got, err := NewRepository(pool).List(context.Background())

	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, "อาหารจานเดียว", got[0].Name)
	require.Equal(t, "เครื่องดื่ม", got[1].Name)
	require.NotEmpty(t, got[0].ID)
}
