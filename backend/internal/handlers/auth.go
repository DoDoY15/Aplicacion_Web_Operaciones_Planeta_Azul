package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/planeta-azul/backend/internal/auth"
	"github.com/planeta-azul/backend/internal/middleware"
	"github.com/planeta-azul/backend/internal/models"
	"github.com/planeta-azul/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store   *repository.MemStore
	authSvc *auth.Service
}

func NewAuthHandler(store *repository.MemStore, authSvc *auth.Service) *AuthHandler {
	return &AuthHandler{store: store, authSvc: authSvc}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "conta desativada"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}

	tokens, err := h.authSvc.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar tokens"})
		return
	}

	// Set httpOnly cookies
	secureCookie := c.Request.TLS != nil
	c.SetCookie("access_token", tokens.AccessToken, int(8*time.Hour/time.Second), "/", "", secureCookie, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(7*24*time.Hour/time.Second), "/auth/refresh", "", secureCookie, true)

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         sanitizeUser(user),
	})
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	tokenStr, _ := c.Cookie("refresh_token")
	if tokenStr == "" {
		var req models.RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token ausente"})
			return
		}
		tokenStr = req.RefreshToken
	}

	userID, err := h.authSvc.ParseUserIDFromRefresh(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido"})
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil || !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não encontrado"})
		return
	}

	tokens, err := h.authSvc.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar tokens"})
		return
	}

	secureCookie := c.Request.TLS != nil
	c.SetCookie("access_token", tokens.AccessToken, int(8*time.Hour/time.Second), "/", "", secureCookie, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(7*24*time.Hour/time.Second), "/auth/refresh", "", secureCookie, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout realizado"})
}

// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, _ := c.Get(middleware.ContextUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id inválido"})
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, sanitizeUser(user))
}

func sanitizeUser(u *models.User) *models.User {
	safe := *u
	safe.PasswordHash = ""
	return &safe
}
