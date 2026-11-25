package menus

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

type fakeRepo struct {
	menus []Menu      // returned by List
	saved []MenuInput // every input that reached Create or Update
	err   error       // returned by List, Update and Delete when set
}

func (f *fakeRepo) List(_ context.Context) ([]Menu, error) {
	return f.menus, f.err
}

func (f *fakeRepo) Create(_ context.Context, in MenuInput) (Menu, error) {
	if f.err != nil {
		return Menu{}, f.err
	}
	f.saved = append(f.saved, in)
	return Menu{ID: uuid.New(), Name: in.Name, Price: *in.Price, CategoryID: in.CategoryID}, nil
}

func (f *fakeRepo) Update(_ context.Context, id uuid.UUID, in MenuInput) (Menu, error) {
	if f.err != nil {
		return Menu{}, f.err
	}
	f.saved = append(f.saved, in)
	return Menu{ID: id, Name: in.Name, Price: *in.Price, CategoryID: in.CategoryID}, nil
}

func (f *fakeRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return f.err
}

func price(s string) *decimal.Decimal {
	d := decimal.RequireFromString(s)
	return &d
}

func TestCreate_TrimsAndSaves(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewUsecase(repo)
	categoryID := uuid.New()

	menu, err := uc.Create(context.Background(), MenuInput{
		Name:       " Pad Kra Pao ",
		Price:      price("59.50"),
		CategoryID: categoryID,
	})

	require.NoError(t, err)
	assert.Equal(t, "Pad Kra Pao", menu.Name)
	require.Len(t, repo.saved, 1)
	assert.Equal(t, categoryID, repo.saved[0].CategoryID, "the repository must receive the same category")
}

func TestCreate_RejectsInvalidInput(t *testing.T) {
	anyCategory := uuid.New()
	tests := []struct {
		name      string
		in        MenuInput
		wantField string
		wantMsg   string
	}{
		{"blank name", MenuInput{Name: "   ", Price: price("10"), CategoryID: anyCategory}, "name", "required"},
		{"blank category", MenuInput{Name: "A", Price: price("10"), CategoryID: uuid.Nil}, "categoryId", "required"},
		{"missing price", MenuInput{Name: "A", Price: nil, CategoryID: anyCategory}, "price", "required"},
		{"negative price", MenuInput{Name: "A", Price: price("-1"), CategoryID: anyCategory}, "price", "must be 0 or more"},
		{"3 decimals", MenuInput{Name: "A", Price: price("10.555"), CategoryID: anyCategory}, "price", "must have at most 2 decimal places"},
		{"too big", MenuInput{Name: "A", Price: price("100000000"), CategoryID: anyCategory}, "price", "must be at most 99999999.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			uc := NewUsecase(repo)

			_, err := uc.Create(context.Background(), tt.in)

			var verr *validation.Error
			require.ErrorAs(t, err, &verr)
			assert.Equal(t, tt.wantMsg, verr.Fields[tt.wantField])
			assert.Empty(t, repo.saved, "invalid input must never reach the database")
		})
	}
}

func TestCreate_AcceptsFreeItemAndTwoDecimals(t *testing.T) {
	for _, p := range []string{"0", "10.5", "10.50", "99999999.99"} {
		uc := NewUsecase(&fakeRepo{})

		_, err := uc.Create(context.Background(), MenuInput{Name: "A", Price: price(p), CategoryID: uuid.New()})

		assert.NoError(t, err, "price %s should be allowed", p)
	}
}

func TestUpdate_PassesNotFoundThrough(t *testing.T) {
	uc := NewUsecase(&fakeRepo{err: ErrNotFound})

	_, err := uc.Update(context.Background(), uuid.New(), MenuInput{Name: "A", Price: price("1"), CategoryID: uuid.New()})

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDelete_PassesNotFoundThrough(t *testing.T) {
	uc := NewUsecase(&fakeRepo{err: ErrNotFound})

	err := uc.Delete(context.Background(), uuid.New())

	assert.ErrorIs(t, err, ErrNotFound)
}
