package packages

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

type fakeRepo struct {
	menus []IncludedMenu
	pkgs  []Package
	saved []PackageInput
	err   error
}

func (f *fakeRepo) List(_ context.Context) ([]Package, error) {
	return f.pkgs, f.err
}

func (f *fakeRepo) GetByID(_ context.Context, id uuid.UUID) (Package, error) {
	for _, p := range f.pkgs {
		if p.ID == id {
			return p, f.err
		}
	}
	return Package{}, ErrNotFound
}

func (f *fakeRepo) Create(_ context.Context, in PackageInput) (Package, error) {
	if f.err != nil {
		return Package{}, f.err
	}
	f.saved = append(f.saved, in)
	return Package{ID: uuid.New(), Name: in.Name, IsActive: *in.IsActive}, nil
}

func (f *fakeRepo) Update(_ context.Context, id uuid.UUID, in PackageInput) (Package, error) {
	if f.err != nil {
		return Package{}, f.err
	}
	f.saved = append(f.saved, in)
	return Package{ID: id, Name: in.Name, IsActive: *in.IsActive}, nil
}

func (f *fakeRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return f.err
}

func (f *fakeRepo) ListMenusUpToLevel(_ context.Context, level int) ([]IncludedMenu, error) {
	var out []IncludedMenu
	for _, m := range f.menus {
		if m.MinPackageLevel <= level {
			out = append(out, m)
		}
	}
	return out, f.err
}

func newTestRouter(repo *fakeRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	validation.Setup()
	r := gin.New()
	NewHandler(NewUsecase(repo)).RegisterRoutes(r.Group("/api/v1"))
	return r
}

func get(r *gin.Engine, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

func send(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

const validBody = `{"name":"  หมูผัก 199 ","level":1,"pricePerPerson":199,"timeLimitMinutes":90,
	"freeBelowHeightCm":100,"overtimeFee":30,"overtimeBlockMinutes":15,"overtimeGraceMinutes":5,
	"wasteFee":50,"wasteUnitGrams":100}`

func TestCreateHandler_Returns201(t *testing.T) {
	repo := &fakeRepo{}

	w := send(newTestRouter(repo), http.MethodPost, "/api/v1/buffet-packages", validBody)

	assert.Equal(t, http.StatusCreated, w.Code)
	require.Len(t, repo.saved, 1, "a valid package must reach the repository")
	assert.Equal(t, "หมูผัก 199", repo.saved[0].Name)
	assert.True(t, *repo.saved[0].IsActive)
}

func TestCreateHandler_Returns400(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantCode string
	}{
		{"empty object", `{}`, "VALIDATION_FAILED"},
		{"blank name", strings.Replace(validBody, `"  หมูผัก 199 "`, `"   "`, 1), "VALIDATION_FAILED"},
		{"zero price", strings.Replace(validBody, `"pricePerPerson":199`, `"pricePerPerson":0`, 1), "VALIDATION_FAILED"},
		{"price is text", strings.Replace(validBody, `"pricePerPerson":199`, `"pricePerPerson":"abc"`, 1), "INVALID_JSON"},
		{"broken JSON", `{"name":`, "INVALID_JSON"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}

			w := send(newTestRouter(repo), http.MethodPost, "/api/v1/buffet-packages", tt.body)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, tt.wantCode, errorCode(t, w))
			assert.Empty(t, repo.saved, "invalid input must never reach the repository")
		})
	}
}

func TestCreateHandler_DuplicateName_Returns409(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{err: ErrNameTaken}), http.MethodPost, "/api/v1/buffet-packages", validBody)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "NAME_TAKEN", errorCode(t, w))
}

func TestUpdateHandler_Returns200(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodPut, "/api/v1/buffet-packages/"+uuid.NewString(), validBody)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateHandler_UnknownID_Returns404(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{err: ErrNotFound}), http.MethodPut,
		"/api/v1/buffet-packages/"+uuid.NewString(), validBody)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "NOT_FOUND", errorCode(t, w))
}

func TestDeleteHandler_Returns204WithNoBody(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodDelete, "/api/v1/buffet-packages/"+uuid.NewString(), "")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func errorCode(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	return body.Error.Code
}

func samplePackage() Package {
	return Package{
		ID:             uuid.New(),
		Name:           "หมูผัก 199",
		Level:          1,
		PricePerPerson: decimal.RequireFromString("199"),
	}
}

func TestListHandler_ReturnsPackages(t *testing.T) {
	pkg := samplePackage()
	w := get(newTestRouter(&fakeRepo{pkgs: []Package{pkg}}), "/api/v1/buffet-packages")

	assert.Equal(t, http.StatusOK, w.Code)
	var body struct {
		Success bool      `json:"success"`
		Data    []Package `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.True(t, body.Success)
	require.Len(t, body.Data, 1)
	assert.Equal(t, "หมูผัก 199", body.Data[0].Name)
	assert.Equal(t, 1, body.Data[0].Level)
}

func TestListHandler_Returns500WithoutDetails(t *testing.T) {
	repo := &fakeRepo{err: errors.New("connection refused: secret-db-host:5432")}
	w := get(newTestRouter(repo), "/api/v1/buffet-packages")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotContains(t, w.Body.String(), "secret-db-host")
}

func TestGetHandler_ReturnsPackage(t *testing.T) {
	pkg := samplePackage()
	w := get(newTestRouter(&fakeRepo{pkgs: []Package{pkg}}), "/api/v1/buffet-packages/"+pkg.ID.String())

	assert.Equal(t, http.StatusOK, w.Code)
	var body struct {
		Data Package `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, pkg.ID, body.Data.ID)
}

func TestGetHandler_UnknownID_Returns404(t *testing.T) {
	w := get(newTestRouter(&fakeRepo{}), "/api/v1/buffet-packages/"+uuid.NewString())

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "NOT_FOUND", errorCode(t, w))
}

func TestGetHandler_InvalidID_Returns404(t *testing.T) {
	w := get(newTestRouter(&fakeRepo{}), "/api/v1/buffet-packages/not-a-uuid")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "NOT_FOUND", errorCode(t, w))
}

func TestListMenusHandler_ReturnsMenusUpToPackageLevel(t *testing.T) {
	pkg := samplePackage() // level 1
	repo := &fakeRepo{
		pkgs: []Package{pkg},
		menus: []IncludedMenu{
			{ID: uuid.New(), Name: "หมูสามชั้น", MinPackageLevel: 1},
			{ID: uuid.New(), Name: "เนื้อวากิว", MinPackageLevel: 2},
		},
	}

	w := get(newTestRouter(repo), "/api/v1/buffet-packages/"+pkg.ID.String()+"/menus")

	assert.Equal(t, http.StatusOK, w.Code)
	var body struct {
		Data []IncludedMenu `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Data, 1, "a level 1 package must not include a level 2 menu")
	assert.Equal(t, "หมูสามชั้น", body.Data[0].Name)
}

func TestListMenusHandler_UnknownPackage_Returns404(t *testing.T) {
	w := get(newTestRouter(&fakeRepo{}), "/api/v1/buffet-packages/"+uuid.NewString()+"/menus")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "NOT_FOUND", errorCode(t, w))
}
