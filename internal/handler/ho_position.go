package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoPositions GET /api/ho-position/list?ho=<id>
func (h *Handler) ListHoPositions(c *gin.Context) {
	hoID := c.Query("ho")
	items, err := h.services.HoPosition.ListByHo(c.Request.Context(), hoID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateHoPosition POST /api/ho-position/create
// Вызывает SQL-процедуру add_ho_position, которая пересчитывает итоговую сумму ХО.
func (h *Handler) CreateHoPosition(c *gin.Context) {
	var input schema.CreateHoPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoPosition.Create(c.Request.Context(),
		input.HoID, input.EmobileID, input.Quantity, input.UnitPrice, input.Note)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateHoPosition PUT /api/ho-position/update:id
func (h *Handler) UpdateHoPosition(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateHoPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.HoPosition.Update(c.Request.Context(), id, input.Quantity, input.UnitPrice, input.Note); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteHoPosition DELETE /api/ho-position/delete:id
func (h *Handler) DeleteHoPosition(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoPosition.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
