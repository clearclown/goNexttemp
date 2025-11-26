package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ablaze/gonexttemp-backend/internal/auth"
	"github.com/ablaze/gonexttemp-backend/internal/config"
	"github.com/ablaze/gonexttemp-backend/internal/handler"
	"github.com/ablaze/gonexttemp-backend/internal/middleware"
	"github.com/ablaze/gonexttemp-backend/internal/repository"
	"github.com/ablaze/gonexttemp-backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Setup logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTAccessExpiry,
		cfg.JWTRefreshExpiry,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	router := setupRouter(cfg, jwtManager, healthHandler, authHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("Starting server", "addr", addr)
	if err := router.Run(addr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func connectDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func setupRouter(
	cfg *config.Config,
	jwtManager *auth.JWTManager,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.CORSOrigins))

	// Health check
	router.GET("/health", healthHandler.Health)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			protected.GET("/auth/me", authHandler.Me)
		}
	}

	return router
}
