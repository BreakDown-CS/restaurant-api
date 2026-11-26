package packages

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/BreakDown-CS/restaurant-api/internal/response"
)

type Handler struct {
	uc *Usecase
}

func NewHandler(uc *Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/buffet-packages", h.list)
	rg.GET("/buffet-packages/:id", h.get)
}

func (h *Handler) list(c *gin.Context) {
	pkgs, err := h.uc.List(c.Request.Context())
	if err != nil {
		h.fail(c, "list buffet packages", err)
		return
	}

	response.OK(c, pkgs)
}

func (h *Handler) get(c *gin.Context) {
	id, ok := packageID(c)
	if !ok {
		return
	}

	pkg, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		h.fail(c, "get buffet package", err)
		return
	}

	response.OK(c, pkg)
}

func packageID(c *gin.Context) (id uuid.UUID, ok bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Buffet package not found", nil)
		return uuid.Nil, false
	}

	return id, true
}

func (h *Handler) fail(c *gin.Context, op string, err error) {
	if errors.Is(err, ErrNotFound) {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Buffet package not found", nil)
		return
	}

	slog.Error(op, "err", err)
	response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
}
