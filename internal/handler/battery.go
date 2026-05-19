package handler

import (
	"net/http"
	"MISPRIS/internal/schema"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBattery(c *gin.Context) {
	var input schema.CreateBatteryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	id, err := h.services.Battery.Create(c.Request.Context(), input.Name, input.Type, input.Capacity, input.Info)
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusCreated, idResponse{ID: id})
}

func (h *Handler) GetBattery(c *gin.Context) {
	id := c.Param("id")
	b, err := h.services.Battery.GetByID(c.Request.Context(), id)
	if err != nil { newErrorResponse(c, http.StatusNotFound, err.Error()); return }
	c.JSON(http.StatusOK, b)
}

func (h *Handler) ListBatteries(c *gin.Context) {
	list, err := h.services.Battery.List(c.Request.Context())
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateBattery(c *gin.Context) {
	id := c.Param("id")
	var input schema.CreateBatteryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	if err := h.services.Battery.Update(c.Request.Context(), id, input.Name, input.Type, input.Capacity, input.Info); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteBattery(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Battery.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
