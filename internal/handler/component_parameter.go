package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListComponentParameters GET /api/component-parameter/list
func (h *Handler) ListComponentParameters(c *gin.Context) {
	items, err := h.services.ComponentParameter.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetComponentParameter GET /api/component-parameter/getComponentParameter:id
func (h *Handler) GetComponentParameter(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.ComponentParameter.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateComponentParameter POST /api/component-parameter/create
func (h *Handler) CreateComponentParameter(c *gin.Context) {
	var input schema.CreateComponentParameterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.ComponentParameter.Create(
		c.Request.Context(),
		input.ComponentType, input.ParameterID, input.OrderNum,
		input.MinVal, input.MaxVal,
	)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateComponentParameter PUT /api/component-parameter/update:id
func (h *Handler) UpdateComponentParameter(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateComponentParameterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.ComponentParameter.Update(
		c.Request.Context(), id, input.OrderNum, input.MinVal, input.MaxVal,
	); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteComponentParameter DELETE /api/component-parameter/delete:id
func (h *Handler) DeleteComponentParameter(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.ComponentParameter.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// GetComponentParametersByType GET /api/component-parameter/byType:type
// Вызывает SQL-функцию get_component_parameters — полные данные о параметрах типа.
func (h *Handler) GetComponentParametersByType(c *gin.Context) {
	componentType := c.Param("type")
	items, err := h.services.ComponentParameter.GetByType(c.Request.Context(), componentType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CopyComponentParameters POST /api/component-parameter/copyFromType
// Вызывает SQL-процедуру copy_component_parameters — наследование параметров.
func (h *Handler) CopyComponentParameters(c *gin.Context) {
	var input schema.CopyComponentParametersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.ComponentParameter.CopyFromType(
		c.Request.Context(), input.FromType, input.ToType,
	); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "parameters copied"})
}
