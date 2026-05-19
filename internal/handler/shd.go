package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListShd GET /api/shd/list
func (h *Handler) ListShd(c *gin.Context) {
	items, err := h.services.Shd.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetShd GET /api/shd/getShd:id
func (h *Handler) GetShd(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.Shd.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateShd POST /api/shd/create
func (h *Handler) CreateShd(c *gin.Context) {
	var input schema.CreateShdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Shd.Create(c.Request.Context(), input.Name, input.ShdType, input.INN, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateShd PUT /api/shd/update:id
func (h *Handler) UpdateShd(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateShdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Shd.Update(c.Request.Context(), id, input.Name, input.ShdType, input.INN, input.Description); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteShd DELETE /api/shd/delete:id
func (h *Handler) DeleteShd(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Shd.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
