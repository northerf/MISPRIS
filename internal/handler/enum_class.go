package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListEnumClasses GET /api/enum-class/list
func (h *Handler) ListEnumClasses(c *gin.Context) {
	items, err := h.services.EnumClass.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetEnumClass GET /api/enum-class/getEnumClass:id
func (h *Handler) GetEnumClass(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.EnumClass.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateEnumClass POST /api/enum-class/create
func (h *Handler) CreateEnumClass(c *gin.Context) {
	var input schema.CreateEnumClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.EnumClass.Create(c.Request.Context(), input.Name, input.ComponentType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateEnumClass PUT /api/enum-class/update:id
func (h *Handler) UpdateEnumClass(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateEnumClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.EnumClass.Update(c.Request.Context(), id, input.Name, input.ComponentType); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteEnumClass DELETE /api/enum-class/delete:id
func (h *Handler) DeleteEnumClass(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.EnumClass.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// GetEnumClassValues GET /api/enum-class/values:id
// Вызывает SQL-функцию get_enum_values — возвращает позиции в порядке order_num.
func (h *Handler) GetEnumClassValues(c *gin.Context) {
	id := c.Param("id")
	values, err := h.services.EnumClass.GetValues(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, values)
}

// ValidateEnumValue POST /api/enum-class/validate
// Проверяет, существует ли значение в перечислении.
func (h *Handler) ValidateEnumValue(c *gin.Context) {
	var input schema.ValidateEnumValueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	valid, err := h.services.EnumClass.ValidateValue(c.Request.Context(), input.EnumClassID, input.Value)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"valid": valid})
}
