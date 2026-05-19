package handler

import (
	"net/http"
	"MISPRIS/internal/schema"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBody(c *gin.Context) {
	var input schema.CreateBodyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	id, err := h.services.Body.Create(c.Request.Context(), input.CarcassID, input.DoorsID, input.WingsID)
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusCreated, idResponse{ID: id})
}
func (h *Handler) GetBody(c *gin.Context) {
	id := c.Param("id")
	b, err := h.services.Body.GetByID(c.Request.Context(), id)
	if err != nil { newErrorResponse(c, http.StatusNotFound, err.Error()); return }
	c.JSON(http.StatusOK, b)
}
func (h *Handler) ListBodies(c *gin.Context) {
	list, err := h.services.Body.List(c.Request.Context())
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusOK, list)
}
func (h *Handler) UpdateBody(c *gin.Context) {
	id := c.Param("id")
	var input schema.CreateBodyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	if err := h.services.Body.Update(c.Request.Context(), id, input.CarcassID, input.DoorsID, input.WingsID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
func (h *Handler) DeleteBody(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Body.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
func (h *Handler) CreateCarcass(c *gin.Context) {
	var input schema.CreateCarcassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	id, err := h.services.Carcass.Create(c.Request.Context(), input.Name, input.Info)
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusCreated, idResponse{ID: id})
}
func (h *Handler) GetCarcass(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.Carcass.GetByID(c.Request.Context(), id)
	if err != nil { newErrorResponse(c, http.StatusNotFound, err.Error()); return }
	c.JSON(http.StatusOK, item)
}
func (h *Handler) ListCarcasses(c *gin.Context) {
	list, err := h.services.Carcass.List(c.Request.Context())
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusOK, list)
}
func (h *Handler) UpdateCarcass(c *gin.Context) {
	id := c.Param("id")
	var input schema.CreateCarcassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	if err := h.services.Carcass.Update(c.Request.Context(), id, input.Name, input.Info); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
func (h *Handler) DeleteCarcass(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Carcass.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
func (h *Handler) CreateDoors(c *gin.Context) {
	var input schema.CreateDoorsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	id, err := h.services.Doors.Create(c.Request.Context(), input.Name, input.Info)
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusCreated, idResponse{ID: id})
}
func (h *Handler) GetDoors(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.Doors.GetByID(c.Request.Context(), id)
	if err != nil { newErrorResponse(c, http.StatusNotFound, err.Error()); return }
	c.JSON(http.StatusOK, item)
}
func (h *Handler) ListDoors(c *gin.Context) {
	list, err := h.services.Doors.List(c.Request.Context())
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusOK, list)
}
func (h *Handler) UpdateDoors(c *gin.Context) {
	id := c.Param("id")
	var input schema.CreateDoorsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	if err := h.services.Doors.Update(c.Request.Context(), id, input.Name, input.Info); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
func (h *Handler) DeleteDoors(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Doors.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
func (h *Handler) CreateWings(c *gin.Context) {
	var input schema.CreateWingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	id, err := h.services.Wings.Create(c.Request.Context(), input.Name, input.Info)
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusCreated, idResponse{ID: id})
}
func (h *Handler) GetWings(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.Wings.GetByID(c.Request.Context(), id)
	if err != nil { newErrorResponse(c, http.StatusNotFound, err.Error()); return }
	c.JSON(http.StatusOK, item)
}
func (h *Handler) ListWings(c *gin.Context) {
	list, err := h.services.Wings.List(c.Request.Context())
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusOK, list)
}
func (h *Handler) UpdateWings(c *gin.Context) {
	id := c.Param("id")
	var input schema.CreateWingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	if err := h.services.Wings.Update(c.Request.Context(), id, input.Name, input.Info); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
func (h *Handler) DeleteWings(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.Wings.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
