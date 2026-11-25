package categories

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	categories []Category
	err        error
}

func (f *fakeRepo) List(_ context.Context) ([]Category, error) {
	return f.categories, f.err
}

func (f *fakeRepo) Create(_ context.Context, in CategoryInput) (Category, error) {
	return Category{ID: uuid.New(), Name: in.Name}, f.err
}

func (f *fakeRepo) Update(_ context.Context, id uuid.UUID, in CategoryInput) (Category, error) {
	return Category{ID: id, Name: in.Name}, f.err
}

func (f *fakeRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return f.err
}

func newTestRouter(repo *fakeRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(NewUsecase(repo)).RegisterRoutes(r.Group("/api/v1"))
	return r
}

func send(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
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

func TestList_ReturnsCategories(t *testing.T) {
	repo := &fakeRepo{categories: []Category{{ID: uuid.New(), Name: "เครื่องดื่ม"}}}
	w := send(newTestRouter(repo), http.MethodGet, "/api/v1/categories", "")

	assert.Equal(t, http.StatusOK, w.Code)
	var body struct {
		Success bool       `json:"success"`
		Data    []Category `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.True(t, body.Success)
	require.Len(t, body.Data, 1)
	assert.Equal(t, "เครื่องดื่ม", body.Data[0].Name)
}

func TestList_Returns500WithoutDetails(t *testing.T) {
	repo := &fakeRepo{err: errors.New("connection refused: secret-db-host:5432")}
	w := send(newTestRouter(repo), http.MethodGet, "/api/v1/categories", "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotContains(t, w.Body.String(), "secret-db-host")
}

func TestCreate_ReturnsCreatedWithTrimmedName(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodPost, "/api/v1/categories", `{"name":" ของหวาน "}`)

	assert.Equal(t, http.StatusCreated, w.Code)
	var body struct {
		Data Category `json:"data"`
	}

	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ของหวาน", body.Data.Name)
}

func TestCreate_Returns400WhenNameBlank(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodPost, "/api/v1/categories", `{"name":"   "}`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "VALIDATION_FAILED", errorCode(t, w))
}

func TestCreate_Returns409WhenNameTaken(t *testing.T) {
	repo := &fakeRepo{err: ErrNameTaken}

	w := send(newTestRouter(repo), http.MethodPost, "/api/v1/categories", `{"name":"เครื่องดื่ม"}`)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "NAME_TAKEN", errorCode(t, w))
}

func TestUpdate_Returns404WhenMissing(t *testing.T) {
	repo := &fakeRepo{err: ErrNotFound}

	w := send(newTestRouter(repo), http.MethodPut, "/api/v1/categories/"+uuid.NewString(), `{"name":"ของหวาน"}`)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDelete_Returns204(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodDelete, "/api/v1/categories/"+uuid.NewString(), "")

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDelete_Returns409WhenInUse(t *testing.T) {
	repo := &fakeRepo{err: ErrInUse}

	w := send(newTestRouter(repo), http.MethodDelete, "/api/v1/categories/"+uuid.NewString(), "")

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "CATEGORY_IN_USE", errorCode(t, w))
}
