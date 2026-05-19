package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoRoles GET /api/ho-role/list
func (h *Handler) ListHoRoles(c *gin.Context) {
	items, err := h.services.HoRole.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoRole GET /api/ho-role/getHoRole:id
func (h *Handler) GetHoRole(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.HoRole.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateHoRole POST /api/ho-role/create
func (h *Handler) CreateHoRole(c *gin.Context) {
	var input schema.CreateHoRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoRole.Create(c.Request.Context(), input.Name, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateHoRole PUT /api/ho-role/update:id
func (h *Handler) UpdateHoRole(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateHoRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.HoRole.Update(c.Request.Context(), id, input.Name, input.Description); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteHoRole DELETE /api/ho-role/delete:id
func (h *Handler) DeleteHoRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoRole.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
