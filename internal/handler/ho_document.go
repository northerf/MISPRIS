package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoDocuments GET /api/ho-document/list?ho=<id>
func (h *Handler) ListHoDocuments(c *gin.Context) {
	hoID := c.Query("ho")
	items, err := h.services.HoDocument.ListByHo(c.Request.Context(), hoID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateHoDocument POST /api/ho-document/create
func (h *Handler) CreateHoDocument(c *gin.Context) {
	var input schema.CreateHoDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoDocument.Create(c.Request.Context(),
		input.HoID, input.DocClassID, input.DocNumber, input.DocDate, input.Note)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteHoDocument DELETE /api/ho-document/delete:id
func (h *Handler) DeleteHoDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoDocument.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
