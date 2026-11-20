package menus

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BreakDown-CS/restaurant-api/internal/response"
	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

// newTestRouter wires handler -> real usecase -> fake repository, so a test
// exercises everything except PostgreSQL.
func newTestRouter(repo *fakeRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	validation.Setup()
	decimal.MarshalJSONWithoutQuotes = true // main.go sets this for the real server

	r := gin.New()
	NewHandler(NewUsecase(repo)).RegisterRoutes(r.Group("/api/v1"))
	return r
}

// send makes an HTTP request without starting a server and returns the recorded reply.
func send(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// decode reads the standard envelope; Data stays raw so each test decodes what it expects.
func decode(t *testing.T, w *httptest.ResponseRecorder) (body struct {
	Success bool                `json:"success"`
	Data    json.RawMessage     `json:"data"`
	Error   *response.ErrorInfo `json:"error"`
},
) {
	t.Helper()
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body), "body: %s", w.Body.String())
	return body
}

func TestList_ReturnsMenus(t *testing.T) {
	repo := &fakeRepo{menus: []Menu{{ID: uuid.New(), Name: "Pad Kra Pao", Price: decimal.RequireFromString("59.5")}}}

	w := send(newTestRouter(repo), http.MethodGet, "/api/v1/menus", "")

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"price":59.5`, "price must be a JSON number, not a string")
}

func TestCreate_Returns201(t *testing.T) {
	repo := &fakeRepo{}

	w := send(newTestRouter(repo), http.MethodPost, "/api/v1/menus",
		`{"name":"  Pad Kra Pao ","price":59.5,"category":"Rice","isAvailable":true}`)

	assert.Equal(t, http.StatusCreated, w.Code)
	body := decode(t, w)
	assert.True(t, body.Success)
	require.Len(t, repo.saved, 1, "a valid menu must reach the repository")
	assert.Equal(t, "Pad Kra Pao", repo.saved[0].Name)
}

func TestCreate_Returns400(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantCode string
		wantKeys []string // fields expected in error.details
	}{
		{"empty object", `{}`, "VALIDATION_FAILED", []string{"name", "price", "category"}},
		{"blank name (usecase rule)", `{"name":"   ","price":10,"category":"Rice"}`, "VALIDATION_FAILED", []string{"name"}},
		{"3 decimals (usecase rule)", `{"name":"A","price":10.555,"category":"Rice"}`, "VALIDATION_FAILED", []string{"price"}},
		{"broken JSON", `{"name":`, "INVALID_JSON", nil},
		{"price is text", `{"name":"A","price":"abc","category":"Rice"}`, "INVALID_JSON", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}

			w := send(newTestRouter(repo), http.MethodPost, "/api/v1/menus", tt.body)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			body := decode(t, w)
			assert.False(t, body.Success)
			require.NotNil(t, body.Error)
			assert.Equal(t, tt.wantCode, body.Error.Code)
			for _, key := range tt.wantKeys {
				assert.Contains(t, body.Error.Details, key)
			}
			assert.Empty(t, repo.saved)
		})
	}
}

func TestUpdate_Returns200(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodPut, "/api/v1/menus/"+uuid.NewString(),
		`{"name":"A","price":1,"category":"Rice"}`)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDelete_Returns204WithNoBody(t *testing.T) {
	w := send(newTestRouter(&fakeRepo{}), http.MethodDelete, "/api/v1/menus/"+uuid.NewString(), "")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestNotFound_Returns404(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		repo   *fakeRepo
	}{
		{"update: id not a UUID", http.MethodPut, "/api/v1/menus/abc", `{"name":"A","price":1,"category":"R"}`, &fakeRepo{}},
		{"update: no such menu", http.MethodPut, "/api/v1/menus/" + uuid.NewString(), `{"name":"A","price":1,"category":"R"}`, &fakeRepo{err: ErrNotFound}},
		{"delete: id not a UUID", http.MethodDelete, "/api/v1/menus/abc", "", &fakeRepo{}},
		{"delete: no such menu", http.MethodDelete, "/api/v1/menus/" + uuid.NewString(), "", &fakeRepo{err: ErrNotFound}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := send(newTestRouter(tt.repo), tt.method, tt.path, tt.body)

			assert.Equal(t, http.StatusNotFound, w.Code)
			assert.Equal(t, "NOT_FOUND", decode(t, w).Error.Code)
		})
	}
}

func TestUnexpectedError_Returns500WithoutDetails(t *testing.T) {
	repo := &fakeRepo{err: errors.New("connection refused: secret-db-host:5432")}

	w := send(newTestRouter(repo), http.MethodGet, "/api/v1/menus", "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "INTERNAL_ERROR", decode(t, w).Error.Code)
	assert.NotContains(t, w.Body.String(), "secret-db-host", "internal errors must never reach the client")
}
