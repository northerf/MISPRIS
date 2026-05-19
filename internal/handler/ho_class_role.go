package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoClassRoles GET /api/ho-class-role/list?ho_class=<id>
func (h *Handler) ListHoClassRoles(c *gin.Context) {
	hoClassID := c.Query("ho_class")
	items, err := h.services.HoClassRole.List(c.Request.Context(), hoClassID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateHoClassRole POST /api/ho-class-role/create
func (h *Handler) CreateHoClassRole(c *gin.Context) {
	var input schema.CreateHoClassRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoClassRole.Create(c.Request.Context(), input.HoClassID, input.HoRoleID, input.IsRequired)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteHoClassRole DELETE /api/ho-class-role/delete:id
func (h *Handler) DeleteHoClassRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoClassRole.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
