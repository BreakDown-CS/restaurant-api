package menus

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
	rg.GET("/menus", h.list)
	rg.POST("/menus", h.create)
	rg.PUT("/menus/:id", h.update)
	rg.DELETE("/menus/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	menus, err := h.uc.List(c.Request.Context())
	if err != nil {
		slog.Error("list menus", "err", err)
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
		return
	}

	response.OK(c, menus)
}

func (h *Handler) create(c *gin.Context) {
	var in MenuInput
	if !bindInput(c, &in) {
		return
	}

	menu, err := h.uc.Create(c.Request.Context(), in)
	if err != nil {
		h.fail(c, "create menu", err)
		return
	}

	response.Created(c, menu)
}

func (h *Handler) update(c *gin.Context) {
	id, ok := menuID(c)
	if !ok {
		return
	}

	var in MenuInput
	if !bindInput(c, &in) {
		return
	}

	menu, err := h.uc.Update(c.Request.Context(), id, in)
	if err != nil {
		h.fail(c, "update menu", err)
		return
	}

	response.OK(c, menu)
}

func (h *Handler) delete(c *gin.Context) {
	id, ok := menuID(c)
	if !ok {
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, "delete menu", err)
		return
	}

	response.NoContent(c)
}

func menuID(c *gin.Context) (id uuid.UUID, ok bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Menu not found", nil)
		return uuid.Nil, false
	}
	return id, true
}

func bindInput(c *gin.Context, in *MenuInput) (ok bool) {
	if err := c.ShouldBindJSON(in); err != nil {
		if fields, ok := validation.Details(err); ok {
			response.Fail(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid input", fields)
			return false
		}
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
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Menu not found", nil)
		return
	}

	slog.Error(op, "err", err)
	response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
}
