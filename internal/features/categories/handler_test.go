package categories

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func newTestRouter(repo *fakeRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(NewUsecase(repo)).RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestList_ReturnsCategories(t *testing.T) {
	repo := &fakeRepo{categories: []Category{{ID: uuid.New(), Name: "เครื่องดื่ม"}}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)

	newTestRouter(repo).ServeHTTP(w, req)

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
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)

	newTestRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotContains(t, w.Body.String(), "secret-db-host")
}
