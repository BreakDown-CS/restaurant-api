package categories

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/BreakDown-CS/restaurant-api/internal/response"
	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

type Handler struct {
	uc *Usecase
}

func NewHandler(uc *Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/categories", h.list)
	rg.POST("/categories", h.create)
	rg.PUT("/categories/:id", h.update)
	rg.DELETE("/categories/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	categories, err := h.uc.List(c.Request.Context())
	if err != nil {
		slog.Error("list categories", "err", err)
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
		return
	}

	response.OK(c, categories)
}

func (h *Handler) create(c *gin.Context) {
	var in CategoryInput
	if !bindInput(c, &in) {
		return
	}

	cat, err := h.uc.Create(c.Request.Context(), in)
	if err != nil {
		h.fail(c, "create category", err)
		return
	}

	response.Created(c, cat)
}

func (h *Handler) update(c *gin.Context) {
	id, ok := categoryID(c)
	if !ok {
		return
	}

	var in CategoryInput
	if !bindInput(c, &in) {
		return
	}

	cat, err := h.uc.Update(c.Request.Context(), id, in)
	if err != nil {
		h.fail(c, "update category", err)
		return
	}

	response.OK(c, cat)
}

func (h *Handler) delete(c *gin.Context) {
	id, ok := categoryID(c)
	if !ok {
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, "delete category", err)
		return
	}

	response.NoContent(c)
}

func categoryID(c *gin.Context) (id uuid.UUID, ok bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
		return uuid.Nil, false
	}

	return id, true
}

func bindInput(c *gin.Context, in *CategoryInput) (ok bool) {
	if err := c.ShouldBindJSON(in); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_JSON", "Request body must be valid JSON", nil)
		return false
	}

	return true
}

func (h *Handler) fail(c *gin.Context, op string, err error) {
	var verr *validation.Error
	if errors.As(err, &verr) {
		response.Fail(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid input", verr.Fields)
		return
	}

	if errors.Is(err, ErrNotFound) {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
		return
	}

	if errors.Is(err, ErrNameTaken) {
		response.Fail(c, http.StatusConflict, "NAME_TAKEN", "Category name already exists", nil)
		return
	}

	if errors.Is(err, ErrInUse) {
		response.Fail(c, http.StatusConflict, "CATEGORY_IN_USE", "Category still has menus", nil)
		return
	}

	slog.Error(op, "err", err)
	response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
}
