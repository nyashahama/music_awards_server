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
	"github.com/joho/godotenv"
	"github.com/nyashahama/music-awards/internal/config"
	"github.com/nyashahama/music-awards/internal/handlers"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	// is git working correctly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
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
	m, err := migrate.New(
		"file://migrations",
		dbCfg.DatabaseURL(),
	)
	if err != nil {
		log.Fatalf("Could not initialize migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run migrations: %v", err)
	}
	log.Println("✅ Migrations applied")

	// 4) Configure GORM with connection pool
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			QueryFields:            true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create GORM instance: %v", err)
	}

	// 5) Initialize services and handlers
	userRepo := repositories.NewUserRepository(gormDB)
	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)

	// Initialize category dependencies
	categoryRepo := repositories.NewCategoryRepository(gormDB)
	categorySvc := services.NewCategoryService(categoryRepo)
	categoryH := handlers.NewCategoryHandler(categorySvc)

	// Initialize nominee dependencies
	nomineeRepo := repositories.NewNomineeRepository(gormDB)
	nomineeCategoryRepo := repositories.NewNomineeCategoryRepository(gormDB)
	nomineeSvc := services.NewNomineeService(
		nomineeRepo,
		categoryRepo,
		nomineeCategoryRepo,
	)
	nomineeH := handlers.NewNomineeHandler(nomineeSvc)

	// Initialize nominee-category dependencies
	nomineeCategorySvc := services.NewNomineeCategoryService(nomineeCategoryRepo)
	nomineeCategoryH := handlers.NewNomineeCategoryHandler(nomineeCategorySvc)

	// Initialize vote dependencies
	voteRepo := repositories.NewVoteRepository(gormDB)
	voteSvc := services.NewVotingMechanismService(voteRepo, userRepo)
	voteH := handlers.NewVoteHandler(voteSvc)

	// 6) Configure Gin router with production settings
	router := gin.New()

	// Production-friendly middleware stack
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	// API routes
	api := router.Group("/api")
	{
		// Authentication
		api.POST("/register", userH.Register)
		api.POST("/login", userH.Login)

		// Public Category APIs
		api.GET("/categories", categoryH.ListCategories)
		api.GET("/categories/active", categoryH.ListActiveCategories)
		api.GET("/categories/:categoryId", categoryH.GetCategory)

		// Public Nominee APIs
		api.GET("/nominees", nomineeH.GetAllNominees)
		api.GET("/nominees/:id", nomineeH.GetNomineeDetails)
	}

	// Protected routes (require authentication)
	protected := router.Group("/api", middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/profile/:id", userH.GetProfile)
		protected.GET("/profile/users", userH.ListAllUsers)
		protected.PUT("/profile/:id", userH.UpdateProfile)
		protected.DELETE("/profile/:id", userH.DeleteAccount)
		protected.PUT("/profile/:id/promote", userH.PromoteUser)

		// Vote routes
		protected.POST("/votes", voteH.CastVote)
		protected.GET("/votes", voteH.GetUserVotes)
		protected.GET("/votes/available", voteH.GetAvailableVotes)
		protected.PUT("/votes/:id", voteH.ChangeVote)
		protected.DELETE("/votes/:id", voteH.DeleteVote)

		// Nominee-Category routes
		protected.GET("/categories/:categoryId/nominees", nomineeCategoryH.GetNominees)
	}

	// Admin-only routes
	admin := router.Group("/api", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Category Admin APIs
		admin.POST("/categories", categoryH.CreateCategory)
		admin.PUT("/categories/:categoryId", categoryH.UpdateCategory)
		admin.DELETE("/categories/:categoryId", categoryH.DeleteCategory)

		// Nominee Admin APIs
		admin.POST("/nominees", nomineeH.CreateNominee)
		admin.PUT("/nominees/:id", nomineeH.UpdateNominee)
		admin.DELETE("/nominees/:id", nomineeH.DeleteNominee)

		// Nominee-Category Admin APIs
		admin.POST("/nominees/:id/categories", nomineeCategoryH.AddCategory)
		admin.DELETE("/nominees/:id/categories/:categoryId", nomineeCategoryH.RemoveCategory)
		admin.PUT("/nominees/:id/categories", nomineeCategoryH.SetCategories)
		admin.GET("/nominees/:id/categories", nomineeCategoryH.GetCategories)

		// Vote Admin APIs
		admin.GET("/votes/category/:category_id", voteH.GetCategoryVotes)
		admin.GET("/votes/all", voteH.GetAllVotes)
	}

	// 7) Configure server with proper timeouts
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 8) Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
