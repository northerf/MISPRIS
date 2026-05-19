package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoActors GET /api/ho-actor/list?ho=<id>
func (h *Handler) ListHoActors(c *gin.Context) {
	hoID := c.Query("ho")
	items, err := h.services.HoActor.ListByHo(c.Request.Context(), hoID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateHoActor POST /api/ho-actor/create
// Вызывает SQL-процедуру set_ho_actor для назначения актора с валидацией.
func (h *Handler) CreateHoActor(c *gin.Context) {
	var input schema.CreateHoActorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoActor.Create(c.Request.Context(), input.HoID, input.HoRoleID, input.ShdID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteHoActor DELETE /api/ho-actor/delete:id
func (h *Handler) DeleteHoActor(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoActor.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
