package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nyashahama/music-awards/internal/config"
	"github.com/nyashahama/music-awards/internal/db"
	"github.com/nyashahama/music-awards/internal/handlers"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
  // 1) Load config
  dbCfg, err := config.LoadDBConfig()
  if err != nil {
    log.Fatalf("Failed to load DB config: %v", err)
  }

  // 2) Open raw *sql.DB
  sqlDB, err := config.InitDB(dbCfg)
  if err != nil {
    log.Fatalf("Failed to init DB: %v", err)
  }
  defer sqlDB.Close()

  // 3) Run file-based migrations
  //    “file://migrations” means “look in <cwd>/migrations”
  m, err := migrate.New(
    "file://migrations",
    dbCfg.DatabaseURL(), // e.g. "postgres://user:pass@host:port/db?sslmode=disable"
  )
  if err != nil {
    log.Fatalf("Could not initialize migrations: %v", err)
  }
  if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatalf("Could not run migrations: %v", err)
  }
  log.Println("✅ Migrations applied")

  // 4) Open GORM on the same *sql.DB
  gormDB, err := gorm.Open(
    postgres.New(postgres.Config{Conn: sqlDB}),
    &gorm.Config{},
  )
  if err != nil {
    log.Fatalf("Failed to create GORM instance: %v", err)
  }

  // 5) (Optional) If you still want GORM AutoMigrate for non-critical models:
  if err := db.MigrateModels(gormDB); err != nil {
    log.Fatalf("GORM model migration failed: %v", err)
  }

  // … the rest of your service initialization …
  userRepo := repositories.NewUserRepository(gormDB.Statement.DB)
  userSvc  := services.NewUserService(userRepo)
  userH    := handlers.NewUserHandler(userSvc)

  router := gin.Default()

	// Add CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // For production, specify your frontend domains instead of "*"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// [Rest of your existing routes setup...]
	api := router.Group("/api")
	{
		api.POST("/register", userH.Register)
		api.POST("/login", userH.Login)
	}
  protected := router.Group("/api", middleware.AuthMiddleware())
  {
    protected.GET("/profile",         handlers.ProfileHandler)
    protected.GET("/profile/:id",     userH.GetProfile)
    protected.PUT("/profile/:id",     userH.UpdateProfile)
    protected.DELETE("/profile/:id",  userH.DeleteAccount)
    protected.PUT("/profile/:id/promote", userH.PromoteUser)
  }

  server := &http.Server{ Addr: ":8080", Handler: router }
  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    log.Printf("Starting server on %s", server.Addr)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      log.Fatalf("Server error: %v", err)
    }
  }()

  <-quit
  log.Println("Shutting down…")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  server.Shutdown(ctx)
}
