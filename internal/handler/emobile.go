package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateEmobile(c *gin.Context) {
	var input schema.CreateEmobileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Emobile.Create(c.Request.Context(), input.Name,
		input.PowerPointID, input.BatteryID, input.ChargerSystemID,
		input.ChassisID, input.BodyID, input.ElectronicsID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, idResponse{ID: id})
}

func (h *Handler) GetEmobile(c *gin.Context) {
	id := c.Param("id")
	emobile, err := h.services.Emobile.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, emobile)
}

func (h *Handler) ListEmobiles(c *gin.Context) {
	list, err := h.services.Emobile.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateEmobile(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateEmobileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Emobile.Update(c.Request.Context(), id, input.Name,
		input.PowerPointID, input.BatteryID, input.ChargerSystemID,
		input.ChassisID, input.BodyID, input.ElectronicsID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteEmobile(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Emobile.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
