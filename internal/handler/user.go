package handler

import (
	"net/http"
	"MISPRIS/internal/schema"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var input schema.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	id, err := h.services.User.Create(c.Request.Context(), input.Username, input.Password, input.Role, input.IsActive)
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	c.JSON(http.StatusCreated, idResponse{ID: id})
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.services.User.GetByID(c.Request.Context(), id)
	if err != nil { newErrorResponse(c, http.StatusNotFound, err.Error()); return }
	if user == nil { newErrorResponse(c, http.StatusNotFound, "user not found"); return }
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.services.User.List(c.Request.Context())
	if err != nil { newErrorResponse(c, http.StatusInternalServerError, err.Error()); return }
	for i := range users { users[i].Password = "" }
	c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input schema.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()); return
	}
	var isActive bool
	if input.IsActive != nil {
		isActive = *input.IsActive
	} else {
		existing, err := h.services.User.GetByID(c.Request.Context(), id)
		if err != nil || existing == nil { newErrorResponse(c, http.StatusNotFound, "user not found"); return }
		isActive = existing.IsActive
	}
	if err := h.services.User.Update(c.Request.Context(), id, input.Username, input.Password, input.Role, isActive); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.User.Delete(c.Request.Context(), id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
