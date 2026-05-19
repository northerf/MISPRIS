package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListParameters GET /api/parameter/list
func (h *Handler) ListParameters(c *gin.Context) {
	items, err := h.services.Parameter.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetParameter GET /api/parameter/getParameter:id
func (h *Handler) GetParameter(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.Parameter.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateParameter POST /api/parameter/create
func (h *Handler) CreateParameter(c *gin.Context) {
	var input schema.CreateParameterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Parameter.Create(
		c.Request.Context(),
		input.Designation, input.Name, input.ParamType,
		input.MeasuringUnit, input.EnumClassID,
	)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateParameter PUT /api/parameter/update:id
func (h *Handler) UpdateParameter(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateParameterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Parameter.Update(
		c.Request.Context(),
		id, input.Designation, input.Name, input.ParamType,
		input.MeasuringUnit, input.EnumClassID,
	); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteParameter DELETE /api/parameter/delete:id
func (h *Handler) DeleteParameter(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Parameter.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
