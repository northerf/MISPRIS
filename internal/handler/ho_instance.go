package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoInstances GET /api/ho/list?ho_class=<id>
func (h *Handler) ListHoInstances(c *gin.Context) {
	hoClassID := c.Query("ho_class")
	items, err := h.services.HoInstance.List(c.Request.Context(), hoClassID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoInstance GET /api/ho/getHo:id
func (h *Handler) GetHoInstance(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.HoInstance.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateHoInstance POST /api/ho/create
func (h *Handler) CreateHoInstance(c *gin.Context) {
	var input schema.CreateHoInstanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoInstance.Create(c.Request.Context(),
		input.HoClassID, input.DocNumber, input.DocDate, input.TotalAmount, input.Note)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateHoInstance PUT /api/ho/update:id
func (h *Handler) UpdateHoInstance(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateHoInstanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.HoInstance.Update(c.Request.Context(), id,
		input.Status, input.DocNumber, input.DocDate, input.TotalAmount, input.Note); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteHoInstance DELETE /api/ho/delete:id
func (h *Handler) DeleteHoInstance(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoInstance.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// FindHoByClass GET /api/ho/findByClass:id
// Вызывает SQL-функцию find_ho_by_class — сводная информация об операциях данного типа.
func (h *Handler) FindHoByClass(c *gin.Context) {
	id := c.Param("id")
	items, err := h.services.HoInstance.FindByClass(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}
