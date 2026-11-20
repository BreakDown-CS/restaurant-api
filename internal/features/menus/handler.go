package menus

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

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
	if err := c.ShouldBindJSON(&in); err != nil {
		if fields, ok := validation.Details(err); ok {
			response.Fail(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid input", fields)
			return
		}
		response.Fail(c, http.StatusBadRequest, "INVALID_JSON", "Request body must be valid JSON", nil)
		return
	}

	menu, err := h.uc.Create(c.Request.Context(), in)
	if err != nil {
		h.fail(c, "create menu", err)
		return
	}

	response.Created(c, menu)
}

func (h *Handler) fail(c *gin.Context, op string, err error) {
	var verr *validation.Error
	if errors.As(err, &verr) {
		response.Fail(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid input", verr.Fields)
		return
	}

	slog.Error(op, "err", err)
	response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
}
