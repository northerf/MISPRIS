package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoClassParameters GET /api/ho-class-parameter/list?ho_class=<id>
func (h *Handler) ListHoClassParameters(c *gin.Context) {
	hoClassID := c.Query("ho_class")
	items, err := h.services.HoClassParameter.List(c.Request.Context(), hoClassID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoClassParameter GET /api/ho-class-parameter/getHoClassParameter:id
func (h *Handler) GetHoClassParameter(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.HoClassParameter.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateHoClassParameter POST /api/ho-class-parameter/create
func (h *Handler) CreateHoClassParameter(c *gin.Context) {
	var input schema.CreateHoClassParameterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoClassParameter.Create(c.Request.Context(),
		input.HoClassID, input.ParameterID, input.OrderNum, input.MinVal, input.MaxVal)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateHoClassParameter PUT /api/ho-class-parameter/update:id
func (h *Handler) UpdateHoClassParameter(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateHoClassParameterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.HoClassParameter.Update(c.Request.Context(), id, input.OrderNum, input.MinVal, input.MaxVal); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteHoClassParameter DELETE /api/ho-class-parameter/delete:id
func (h *Handler) DeleteHoClassParameter(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoClassParameter.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// CopyHoClassParameters POST /api/ho-class-parameter/copyFromClass
// Копирует параметры из одного класса ХО в другой.
func (h *Handler) CopyHoClassParameters(c *gin.Context) {
	var input schema.CopyHoClassParametersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.HoClassParameter.CopyFromClass(c.Request.Context(), input.FromClassID, input.ToClassID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "copied"})
}
