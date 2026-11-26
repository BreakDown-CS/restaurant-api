package packages

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

func validInput() PackageInput {
	return PackageInput{
		Name:                 "หมูผัก 199",
		Level:                1,
		PricePerPerson:       decimal.RequireFromString("199"),
		TimeLimitMinutes:     90,
		FreeBelowHeightCm:    100,
		OvertimeFee:          decimal.RequireFromString("30"),
		OvertimeBlockMinutes: 15,
		OvertimeGraceMinutes: 5,
		WasteFee:             decimal.RequireFromString("50"),
		WasteUnitGrams:       100,
	}
}

func TestCreate_TrimsNameAndDefaultsToActive(t *testing.T) {
	repo := &fakeRepo{}
	in := validInput()
	in.Name = "  หมูผัก 199  "

	_, err := NewUsecase(repo).Create(context.Background(), in)

	require.NoError(t, err)
	require.Len(t, repo.saved, 1)
	assert.Equal(t, "หมูผัก 199", repo.saved[0].Name)
	require.NotNil(t, repo.saved[0].IsActive)
	assert.True(t, *repo.saved[0].IsActive, "a missing isActive must become true")
}

func TestCreate_KeepsIsActiveFalse(t *testing.T) {
	repo := &fakeRepo{}
	in := validInput()
	in.IsActive = new(false)

	_, err := NewUsecase(repo).Create(context.Background(), in)

	require.NoError(t, err)
	assert.False(t, *repo.saved[0].IsActive)
}

func TestCreate_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		change    func(in *PackageInput)
		wantField string
		wantMsg   string
	}{
		{"blank name", func(in *PackageInput) { in.Name = "   " }, "name", "required"},
		{"zero price", func(in *PackageInput) { in.PricePerPerson = decimal.Zero }, "pricePerPerson", "must be more than 0"},
		{"3 decimals", func(in *PackageInput) { in.PricePerPerson = decimal.RequireFromString("199.999") }, "pricePerPerson", "must have at most 2 decimal places"},
		{"price too big", func(in *PackageInput) { in.PricePerPerson = decimal.RequireFromString("100000000") }, "pricePerPerson", "must be at most 99999999.99"},
		{"negative overtime fee", func(in *PackageInput) { in.OvertimeFee = decimal.RequireFromString("-1") }, "overtimeFee", "must be 0 or more"},
		{"waste fee 3 decimals", func(in *PackageInput) { in.WasteFee = decimal.RequireFromString("1.005") }, "wasteFee", "must have at most 2 decimal places"},
		{"zero level", func(in *PackageInput) { in.Level = 0 }, "level", "must be more than 0"},
		{"zero time limit", func(in *PackageInput) { in.TimeLimitMinutes = 0 }, "timeLimitMinutes", "must be more than 0"},
		{"zero block", func(in *PackageInput) { in.OvertimeBlockMinutes = 0 }, "overtimeBlockMinutes", "must be more than 0"},
		{"zero waste unit", func(in *PackageInput) { in.WasteUnitGrams = 0 }, "wasteUnitGrams", "must be more than 0"},
		{"negative height", func(in *PackageInput) { in.FreeBelowHeightCm = -1 }, "freeBelowHeightCm", "must be 0 or more"},
		{"negative grace", func(in *PackageInput) { in.OvertimeGraceMinutes = -1 }, "overtimeGraceMinutes", "must be 0 or more"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			in := validInput()
			tt.change(&in)

			_, err := NewUsecase(repo).Create(context.Background(), in)

			var verr *validation.Error
			require.ErrorAs(t, err, &verr)
			assert.Equal(t, tt.wantMsg, verr.Fields[tt.wantField])
			assert.Empty(t, repo.saved, "invalid input must never reach the database")
		})
	}
}

func TestCreate_AcceptsZeroFeesAndZeroHeight(t *testing.T) {
	in := validInput()
	in.OvertimeFee = decimal.Zero
	in.WasteFee = decimal.Zero
	in.FreeBelowHeightCm = 0
	in.OvertimeGraceMinutes = 0

	_, err := NewUsecase(&fakeRepo{}).Create(context.Background(), in)

	assert.NoError(t, err)
}

func TestUpdate_ValidatesAndPassesNotFoundThrough(t *testing.T) {
	bad := validInput()
	bad.Level = 0
	_, err := NewUsecase(&fakeRepo{}).Update(context.Background(), uuid.New(), bad)
	var verr *validation.Error
	require.ErrorAs(t, err, &verr)

	_, err = NewUsecase(&fakeRepo{err: ErrNotFound}).Update(context.Background(), uuid.New(), validInput())
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDelete_PassesNotFoundThrough(t *testing.T) {
	err := NewUsecase(&fakeRepo{err: ErrNotFound}).Delete(context.Background(), uuid.New())

	assert.ErrorIs(t, err, ErrNotFound)
}
