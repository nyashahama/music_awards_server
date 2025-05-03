package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyashahama/music-awards/internal/config"
	"github.com/nyashahama/music-awards/internal/db"
	"github.com/nyashahama/music-awards/internal/handlers"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	// Load config and initialize DB
	dbCfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Failed to load DB config: %v", err)
	}

	sqlDB, err := config.InitDB(dbCfg)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("Failed to create GORM instance: %v", err)
	}

	// Migrate models
	if err := db.MigrateModels(gormDB,
		&models.User{},
		&models.Category{},
		&models.Nominee{},
		&models.Vote{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Instantiate services and handlers
	userRepo := repositories.NewUserRepository(gormDB.Statement.DB)
	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)

	// Initialize Gin router
	router := gin.Default()

	// Public routes
	api := router.Group("/api")
	{
		api.POST("/register", userH.Register)
		api.POST("/login", userH.Login)
	}

	// Protected routes
	protected := router.Group("/api", middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.ProfileHandler)
		protected.GET("/profile/:id", userH.GetProfile)
		protected.PUT("/profile/:id", userH.UpdateProfile)
		protected.DELETE("/profile/:id", userH.DeleteAccount)
		protected.PUT("/profile/:id/promote", userH.PromoteUser)
	}

	// Configure and start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
