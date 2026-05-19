package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListEnumPositions GET /api/enum-position/list
func (h *Handler) ListEnumPositions(c *gin.Context) {
	items, err := h.services.EnumPosition.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetEnumPosition GET /api/enum-position/getEnumPosition:id
func (h *Handler) GetEnumPosition(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.EnumPosition.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateEnumPosition POST /api/enum-position/create
func (h *Handler) CreateEnumPosition(c *gin.Context) {
	var input schema.CreateEnumPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.EnumPosition.Create(
		c.Request.Context(), input.EnumClassID, input.Value, input.OrderNum,
	)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateEnumPosition PUT /api/enum-position/update:id
func (h *Handler) UpdateEnumPosition(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateEnumPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.EnumPosition.Update(
		c.Request.Context(), id, input.Value, input.OrderNum,
	); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteEnumPosition DELETE /api/enum-position/delete:id
func (h *Handler) DeleteEnumPosition(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.EnumPosition.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// ReorderEnumPosition POST /api/enum-position/reorder:id
// FIX: сначала читаем текущее значение, затем обновляем только order_num —
// иначе value затирается пустой строкой.
func (h *Handler) ReorderEnumPosition(c *gin.Context) {
	id := c.Param("id")
	var input schema.ReorderEnumPositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// Загружаем текущую позицию, чтобы сохранить её value
	existing, err := h.services.EnumPosition.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "позиция не найдена")
		return
	}
	if err := h.services.EnumPosition.Update(
		c.Request.Context(), id, existing.Value, input.NewOrderNum,
	); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "reordered", "new_order_num": input.NewOrderNum})
}
