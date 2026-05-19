package handler

import (
	"net/http"

	"MISPRIS/internal/schema"

	"github.com/gin-gonic/gin"
)

// ListDocumentClasses GET /api/document-class/list
func (h *Handler) ListDocumentClasses(c *gin.Context) {
	items, err := h.services.DocumentClass.List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetDocumentClass GET /api/document-class/getDocumentClass:id
func (h *Handler) GetDocumentClass(c *gin.Context) {
	id := c.Param("id")
	item, err := h.services.DocumentClass.GetByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateDocumentClass POST /api/document-class/create
func (h *Handler) CreateDocumentClass(c *gin.Context) {
	var input schema.CreateDocumentClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.DocumentClass.Create(c.Request.Context(), input.Name, input.Code, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateDocumentClass PUT /api/document-class/update:id
func (h *Handler) UpdateDocumentClass(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateDocumentClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.DocumentClass.Update(c.Request.Context(), id, input.Name, input.Code, input.Description); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteDocumentClass DELETE /api/document-class/delete:id
func (h *Handler) DeleteDocumentClass(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.DocumentClass.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
