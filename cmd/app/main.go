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
<<<<<<< Updated upstream
=======
	// Load environment variables
>>>>>>> Stashed changes
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
	log.Println("✅ Migrations applied successfully")

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

<<<<<<< Updated upstream
	// 5) Initialize services and handlers
	userRepo := repositories.NewUserRepository(gormDB)
	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)

	// Initialize category dependencies
=======
	// 5) Initialize email service
	emailCfg := config.LoadEmailConfig()
	emailService := services.NewEmailService(emailCfg)

	// 6) Initialize repositories
	userRepo := repositories.NewUserRepository(gormDB)
>>>>>>> Stashed changes
	categoryRepo := repositories.NewCategoryRepository(gormDB)
	nomineeRepo := repositories.NewNomineeRepository(gormDB)
	nomineeCategoryRepo := repositories.NewNomineeCategoryRepository(gormDB)
	voteRepo := repositories.NewVoteRepository(gormDB)

	// 7) Initialize services
	passwordResetService := services.NewPasswordResetService(userRepo, emailService)
	userSvc := services.NewUserService(userRepo, passwordResetService, emailService)
	categorySvc := services.NewCategoryService(categoryRepo)
	nomineeSvc := services.NewNomineeService(nomineeRepo, categoryRepo, nomineeCategoryRepo)
	nomineeCategorySvc := services.NewNomineeCategoryService(nomineeCategoryRepo)

	// Initialize improved voting service with all dependencies
	votingSvc := services.NewVotingService(
		voteRepo,
		userRepo,
		categoryRepo,
		nomineeRepo,
		nomineeCategoryRepo,
	)

	// 8) Initialize handlers
	userH := handlers.NewUserHandler(userSvc)
	categoryH := handlers.NewCategoryHandler(categorySvc)
	nomineeH := handlers.NewNomineeHandler(nomineeSvc)
	nomineeCategoryH := handlers.NewNomineeCategoryHandler(nomineeCategorySvc)
	voteH := handlers.NewVoteHandler(votingSvc)

	// 9) Configure Gin router with production settings
	router := gin.New()

	// Production-friendly middleware stack
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
<<<<<<< Updated upstream
			AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
=======
			AllowOrigins: []string{
				os.Getenv("FRONTEND_URL"),
				"*",
				"https://music-awards-web.onrender.com",
			},
>>>>>>> Stashed changes
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

<<<<<<< Updated upstream
	// API routes
	api := router.Group("/api")
	{
		// Authentication
		api.POST("/register", userH.Register)
		api.POST("/login", userH.Login)
=======
	// 10) Register routes
	setupRoutes(router, userH, categoryH, nomineeH, nomineeCategoryH, voteH)
>>>>>>> Stashed changes

	// 11) Configure server with proper timeouts
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

	// 12) Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("⏳ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

// setupRoutes configures all API routes
func setupRoutes(
	router *gin.Engine,
	userH *handlers.UserHandler,
	categoryH *handlers.CategoryHandler,
	nomineeH *handlers.NomineeHandler,
	nomineeCategoryH *handlers.NomineeCategoryHandler,
	voteH *handlers.VoteHandler,
) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API v1 group
	api := router.Group("/api")
	{
		// ====================================================================
		// PUBLIC ROUTES (No authentication required)
		// ====================================================================

		// Authentication endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", userH.Register)
			auth.POST("/login", userH.Login)
			auth.POST("/forgot-password", userH.ForgotPassword)
			auth.POST("/reset-password", userH.ResetPassword)
			auth.POST("/validate-reset-token", userH.ValidateResetToken)
		}

		// Public category endpoints
		categories := api.Group("/categories")
		{
			categories.GET("", categoryH.ListCategories)
			categories.GET("/active", categoryH.ListActiveCategories)
			categories.GET("/:categoryId", categoryH.GetCategory)
			// Public endpoint to get nominees for a category
			categories.GET("/:categoryId/nominees", nomineeCategoryH.GetNominees)
		}

		// Public nominee endpoints
		nominees := api.Group("/nominees")
		{
			nominees.GET("", nomineeH.GetAllNominees)
			nominees.GET("/:id", nomineeH.GetNomineeDetails)
		}

		// ====================================================================
		// PROTECTED ROUTES (Authentication required)
		// ====================================================================

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile endpoints
			profile := protected.Group("/profile")
			{
				profile.GET("/:id", userH.GetProfile)
				profile.PUT("/:id", userH.UpdateProfile)
				profile.DELETE("/:id", userH.DeleteAccount)
			}

			// Voting endpoints (updated with new structure)
			votes := protected.Group("/votes")
			{
				votes.POST("", voteH.CastVote)
				votes.GET("/me", voteH.GetMyVotes)
				votes.GET("/me/summary", voteH.GetMyVoteSummary)
				votes.GET("/me/available", voteH.GetAvailableVotes)
				votes.PUT("/:id", voteH.ChangeVote)
				votes.DELETE("/:id", voteH.DeleteVote)
			}
		}

		// ====================================================================
		// ADMIN ROUTES (Authentication + Admin role required)
		// ====================================================================

		admin := api.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// User management
			users := admin.Group("/users")
			{
				users.GET("", userH.ListAllUsers)
				users.POST("/:id/promote", userH.PromoteUser)
			}

			// Category management
			categoryAdmin := admin.Group("/categories")
			{
				categoryAdmin.POST("", categoryH.CreateCategory)
				categoryAdmin.PUT("/:categoryId", categoryH.UpdateCategory)
				categoryAdmin.DELETE("/:categoryId", categoryH.DeleteCategory)
			}

			// Nominee management
			nomineeAdmin := admin.Group("/nominees")
			{
				nomineeAdmin.POST("", nomineeH.CreateNominee)
				nomineeAdmin.PUT("/:id", nomineeH.UpdateNominee)
				nomineeAdmin.DELETE("/:id", nomineeH.DeleteNominee)
			}

			// Nominee-Category relationship management
			nomineeCategoryAdmin := admin.Group("/nominees/:id/categories")
			{
				nomineeCategoryAdmin.POST("", nomineeCategoryH.AddCategory)
				nomineeCategoryAdmin.DELETE("/:categoryId", nomineeCategoryH.RemoveCategory)
				nomineeCategoryAdmin.PUT("", nomineeCategoryH.SetCategories)
				nomineeCategoryAdmin.GET("", nomineeCategoryH.GetCategories)
			}

			// Vote management and analytics
			voteAdmin := admin.Group("/votes")
			{
				voteAdmin.GET("/all", voteH.GetAllVotes)
				// voteAdmin.GET("/category/:category_id", voteH.GetCategoryVotes)
				voteAdmin.GET("/category/:category_id/stats", voteH.GetCategoryStats)
				voteAdmin.GET("/nominee/:nominee_id/stats", voteH.GetNomineeStats)
			}
		}
	}

	// Log registered routes (helpful for debugging)
	if gin.Mode() == gin.DebugMode {
		log.Println("\n📋 Registered Routes:")
		for _, route := range router.Routes() {
			log.Printf("  %s %s", route.Method, route.Path)
		}
	}
}
