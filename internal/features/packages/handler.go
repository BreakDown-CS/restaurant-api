package packages

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
	rg.GET("/buffet-packages", h.list)
	rg.GET("/buffet-packages/:id", h.get)
	rg.POST("/buffet-packages", h.create)
	rg.PUT("/buffet-packages/:id", h.update)
	rg.DELETE("/buffet-packages/:id", h.delete)
	rg.GET("/buffet-packages/:id/menus", h.listMenus)
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

func (h *Handler) listMenus(c *gin.Context) {
	id, ok := packageID(c)
	if !ok {
		return
	}

	menus, err := h.uc.ListMenus(c.Request.Context(), id)
	if err != nil {
		h.fail(c, "list buffet package menus", err)
		return
	}

	response.OK(c, menus)
}

func (h *Handler) create(c *gin.Context) {
	var in PackageInput
	if !bindInput(c, &in) {
		return
	}

	pkg, err := h.uc.Create(c.Request.Context(), in)
	if err != nil {
		h.fail(c, "create buffet package", err)
		return
	}

	response.Created(c, pkg)
}

func (h *Handler) update(c *gin.Context) {
	id, ok := packageID(c)
	if !ok {
		return
	}

	var in PackageInput
	if !bindInput(c, &in) {
		return
	}

	pkg, err := h.uc.Update(c.Request.Context(), id, in)
	if err != nil {
		h.fail(c, "update buffet package", err)
		return
	}

	response.OK(c, pkg)
}

func (h *Handler) delete(c *gin.Context) {
	id, ok := packageID(c)
	if !ok {
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, "delete buffet package", err)
		return
	}

	response.NoContent(c)
}

func packageID(c *gin.Context) (id uuid.UUID, ok bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Buffet package not found", nil)
		return uuid.Nil, false
	}

	return id, true
}

func bindInput(c *gin.Context, in *PackageInput) (ok bool) {
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
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "Buffet package not found", nil)
		return
	}

	if errors.Is(err, ErrNameTaken) {
		response.Fail(c, http.StatusConflict, "NAME_TAKEN", "Buffet package name already exists", nil)
		return
	}

	slog.Error(op, "err", err)
	response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", nil)
}
