package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/planeta-azul/backend/internal/auth"
	"github.com/planeta-azul/backend/internal/config"
	"github.com/planeta-azul/backend/internal/handlers"
	"github.com/planeta-azul/backend/internal/middleware"
	"github.com/planeta-azul/backend/internal/models"
	"github.com/planeta-azul/backend/internal/repository"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Auth service — use key files if available, otherwise generate ephemeral keys (dev/no-DB mode)
	var authSvc *auth.Service
	if _, err := os.Stat(cfg.JWTPrivateKeyPath); err == nil {
		authSvc, err = auth.NewService(cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath, cfg.AccessExpiry, cfg.RefreshExpiry)
		if err != nil {
			log.Fatalf("failed to init auth service: %v", err)
		}
		log.Println("✓ JWT: using RSA key files")
	} else {
		log.Println("⚠ JWT key files not found — generating ephemeral RSA keys (dev mode)")
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("failed to generate RSA key: %v", err)
		}
		authSvc = auth.NewServiceFromKeys(privateKey, &privateKey.PublicKey, 8*time.Hour, 168*time.Hour)
	}

	// In-memory store (no DB required)
	store := repository.NewMemStore()
	log.Println("✓ Store: using in-memory store (no DB)")

	// Handlers
	authHandler := handlers.NewAuthHandler(store, authSvc)
	itemHandler := handlers.NewItemHandler(store)
	userHandler := handlers.NewUserHandler(store)

	// Router
	r := gin.Default()

	allowedOrigins := strings.Split(getenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001"), ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "mode": "in-memory", "env": cfg.Env})
	})

	// Auth routes (public)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
	}

	// Protected routes
	api := r.Group("/api/v1")
	api.Use(middleware.RequireAuth(authSvc))
	{
		// Auth
		api.GET("/auth/me", authHandler.Me)

		// Users
		api.GET("/users", userHandler.List)
		api.GET("/users/:id", userHandler.Get)
		api.GET("/areas", userHandler.ListAreas)
		api.GET("/notifications", userHandler.GetNotifications)

		// Items (Módulo 1)
		items := api.Group("/items")
		{
			items.GET("", itemHandler.List)
			items.POST("", itemHandler.Create)
			items.GET("/:id", itemHandler.Get)
			items.PATCH("/:id", itemHandler.Update)
			items.DELETE("/:id", middleware.RequireMinRole(models.RoleChefArea), itemHandler.Delete)
			items.GET("/:id/comments", itemHandler.ListComments)
			items.POST("/:id/comments", itemHandler.AddComment)
			items.GET("/:id/interactions", itemHandler.ListInteractions)
			items.POST("/:id/interactions", itemHandler.CreateInteraction)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Planeta Azul backend running on %s", addr)
	log.Println("   Seed credentials:")
	log.Println("   admin@planetaazul.com  / admin123")
	log.Println("   chefe@planetaazul.com  / chefe123")
	log.Println("   sup@planetaazul.com    / sup123")
	log.Println("   membro@planetaazul.com / membro123")

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
