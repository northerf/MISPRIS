package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListHoClassDocuments GET /api/ho-class-document/list?ho_class=<id>
func (h *Handler) ListHoClassDocuments(c *gin.Context) {
	hoClassID := c.Query("ho_class")
	items, err := h.services.HoClassDocument.ListByClass(c.Request.Context(), hoClassID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateHoClassDocument POST /api/ho-class-document/create
func (h *Handler) CreateHoClassDocument(c *gin.Context) {
	var input schema.CreateHoClassDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.HoClassDocument.Create(c.Request.Context(), input.HoClassID, input.DocClassID, input.RoleName, input.IsRequired)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteHoClassDocument DELETE /api/ho-class-document/delete:id
func (h *Handler) DeleteHoClassDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.HoClassDocument.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
