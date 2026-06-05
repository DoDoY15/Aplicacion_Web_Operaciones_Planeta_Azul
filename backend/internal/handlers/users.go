package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/planeta-azul/backend/internal/repository"
)

type UserHandler struct {
	store *repository.MemStore
}

func NewUserHandler(store *repository.MemStore) *UserHandler {
	return &UserHandler{store: store}
}

// GET /users
func (h *UserHandler) List(c *gin.Context) {
	users := h.store.ListUsers()
	// Strip password hashes
	for _, u := range users {
		u.PasswordHash = ""
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "total": len(users)})
}

// GET /users/:id
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	user, err := h.store.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// GET /users/areas
func (h *UserHandler) ListAreas(c *gin.Context) {
	areas := h.store.ListAreas()
	c.JSON(http.StatusOK, gin.H{"areas": areas})
}

// GET /notifications
func (h *UserHandler) GetNotifications(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	notifications := h.store.GetNotificationsByUser(userID)
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}
