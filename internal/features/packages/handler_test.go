package packages

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	pkgs []Package
	err  error
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

func newTestRouter(repo *fakeRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(NewUsecase(repo)).RegisterRoutes(r.Group("/api/v1"))
	return r
}

func get(r *gin.Engine, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
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
