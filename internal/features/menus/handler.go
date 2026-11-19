package menus

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BreakDown-CS/restaurant-api/internal/response"
)

type Handler struct {
	uc *Usecase
}

func NewHandler(uc *Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/menus", h.list)
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
