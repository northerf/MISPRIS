package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListEmobileParameterValues GET /api/emobile-parameter/list
func (h *Handler) ListEmobileParameterValues(c *gin.Context) {
	items, err := h.services.EmobileParameterValue.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetEmobileParameterValue GET /api/emobile-parameter/getEmobileParameter:id
func (h *Handler) GetEmobileParameterValue(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.EmobileParameterValue.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateEmobileParameterValue POST /api/emobile-parameter/create
func (h *Handler) CreateEmobileParameterValue(c *gin.Context) {
	var input schema.CreateEmobileParameterValueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.EmobileParameterValue.Create(
		c.Request.Context(),
		input.EmobileID, input.ComponentParameterID,
		input.ValReal, input.ValInt, input.ValStr, input.EnumValID,
	)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateEmobileParameterValue PUT /api/emobile-parameter/update:id
func (h *Handler) UpdateEmobileParameterValue(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateEmobileParameterValueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.EmobileParameterValue.Update(
		c.Request.Context(), id,
		input.ValReal, input.ValInt, input.ValStr, input.EnumValID,
	); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteEmobileParameterValue DELETE /api/emobile-parameter/delete:id
func (h *Handler) DeleteEmobileParameterValue(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.EmobileParameterValue.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// GetEmobileParameterValuesByEmobile GET /api/emobile-parameter/byEmobile:id
// Возвращает все значения параметров конкретного автомобиля.
func (h *Handler) GetEmobileParameterValuesByEmobile(c *gin.Context) {
	id := c.Param("id")
	items, err := h.services.EmobileParameterValue.GetByEmobile(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}
