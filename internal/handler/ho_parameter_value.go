package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoParameterValues GET /api/ho-param-value/list?ho=<id>
func (h *Handler) ListHoParameterValues(c *gin.Context) {
	hoID := c.Query("ho")
	items, err := h.services.HoParameterValue.ListByHo(c.Request.Context(), hoID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoParameterValue GET /api/ho-param-value/getHoParamValue:id
func (h *Handler) GetHoParameterValue(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.HoParameterValue.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateHoParameterValue POST /api/ho-param-value/create
// Вызывает SQL-процедуру write_ho_par с валидацией.
func (h *Handler) CreateHoParameterValue(c *gin.Context) {
	var input schema.CreateHoParameterValueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoParameterValue.Create(c.Request.Context(),
		input.HoID, input.HoClassParameterID, input.ValReal, input.ValInt, input.ValStr, input.ValDate, input.EnumValID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateHoParameterValue PUT /api/ho-param-value/update:id
func (h *Handler) UpdateHoParameterValue(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateHoParameterValueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// Retrieve existing record to pass ho_id and ho_class_parameter_id to write_ho_par
	existing, err := h.services.HoParameterValue.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	if err := h.services.HoParameterValue.Update(c.Request.Context(),
		id, existing.HoID, existing.HoClassParameterID,
		input.ValReal, input.ValInt, input.ValStr, input.ValDate, input.EnumValID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteHoParameterValue DELETE /api/ho-param-value/delete:id
func (h *Handler) DeleteHoParameterValue(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoParameterValue.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
