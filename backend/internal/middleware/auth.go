package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/planeta-azul/backend/internal/auth"
	"github.com/planeta-azul/backend/internal/models"
)

const (
	ContextUserID = "user_id"
	ContextEmail  = "email"
	ContextRole   = "role"
	ContextAreaID = "area_id"
)

// RequireAuth validates the access token from cookie or Authorization header
func RequireAuth(authSvc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token ausente"})
			return
		}

		claims, err := authSvc.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, string(claims.Role))
		c.Set(ContextAreaID, claims.AreaID)
		c.Next()
	}
}

// RequireRoles restricts access to specific roles
func RequireRoles(roles ...models.UserRole) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[string(r)] = true
	}

	return func(c *gin.Context) {
		role, _ := c.Get(ContextRole)
		if !allowed[role.(string)] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permisos insuficientes"})
			return
		}
		c.Next()
	}
}

// RequireMinRole allows a role and everything above it in hierarchy
func RequireMinRole(minRole models.UserRole) gin.HandlerFunc {
	hierarchy := map[models.UserRole]int{
		models.RoleMembro:     1,
		models.RoleSupervisor: 2,
		models.RoleChefArea:   3,
		models.RoleChefGeral:  4,
		models.RoleAdmin:      5,
	}

	minLevel := hierarchy[minRole]

	return func(c *gin.Context) {
		roleStr, _ := c.Get(ContextRole)
		level := hierarchy[models.UserRole(roleStr.(string))]
		if level < minLevel {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permisos insuficientes"})
			return
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	// 1. Cookie httpOnly (preferido)
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		return cookie
	}
	// 2. Authorization header (fallback para API clients)
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
