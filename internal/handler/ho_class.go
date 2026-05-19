package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoClasses GET /api/ho-class/list
func (h *Handler) ListHoClasses(c *gin.Context) {
	items, err := h.services.HoClass.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoClass GET /api/ho-class/getHoClass:id
func (h *Handler) GetHoClass(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.HoClass.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateHoClass POST /api/ho-class/create
func (h *Handler) CreateHoClass(c *gin.Context) {
	var input schema.CreateHoClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoClass.Create(c.Request.Context(), input.Name, input.Designation, input.ParentID, input.IsTerminal)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateHoClass PUT /api/ho-class/update:id
func (h *Handler) UpdateHoClass(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateHoClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.HoClass.Update(c.Request.Context(), id, input.Name, input.Designation, input.ParentID, input.IsTerminal); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteHoClass DELETE /api/ho-class/delete:id
func (h *Handler) DeleteHoClass(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoClass.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// GetHoClassTerminal GET /api/ho-class/terminal
// Возвращает только терминальные (листовые) узлы классификатора.
func (h *Handler) GetHoClassTerminal(c *gin.Context) {
	items, err := h.services.HoClass.GetTerminal(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoClassChildren GET /api/ho-class/children:id
// Возвращает дочерние узлы по parentID.
func (h *Handler) GetHoClassChildren(c *gin.Context) {
	id := c.Param("id")
	items, err := h.services.HoClass.GetChildren(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetHoClassParameters GET /api/ho-class/parameters:id
// Вызывает SQL-функцию get_ho_class_parameters — полные данные параметров класса ХО.
func (h *Handler) GetHoClassParameters(c *gin.Context) {
	id := c.Param("id")
	items, err := h.services.HoClassParameter.GetByHoClass(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}
