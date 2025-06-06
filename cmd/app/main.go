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

	//Initialize category dependencies
	categoryRepo := repositories.NewCategoryRepository(gormDB)
	categorySvc := services.NewCategoryService(categoryRepo)
	categoryH := handlers.NewCategoryHandler(categorySvc)

	//Initialize nominee dependencies
	nomineeRepo := repositories.NewNomineeRepository(gormDB)
	nomineeSvc := services.NewNomineeService(nomineeRepo)
	nomineeH := handlers.NewNomineeHandler(nomineeSvc)

	//Initialize nominee

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
		api.POST("/register", userH.Register)
		api.POST("/login", userH.Login)

		//Public Category APIs
		api.GET("/categories", categoryH.ListCategories)
		api.GET("/categories/active", categoryH.ListActiveCategories)
		api.GET("/categories/:categoryId", categoryH.GetCategory)

		//Public Nominee APIs
		api.GET("/nominees", nomineeH.GetAllNominees)
		api.GET("/nominees/:id", nomineeH.GetNomineeDetails)
	}

	// Protected routes
	protected := router.Group("/api", middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.ProfileHandler)
		protected.GET("/profile/:id", userH.GetProfile)
		protected.GET("/profile/users", userH.ListAllUsers)
		protected.PUT("/profile/:id", userH.UpdateProfile)
		protected.DELETE("/profile/:id", userH.DeleteAccount)
		protected.PUT("/profile/:id/promote", userH.PromoteUser)

		//Protected Category APIs
		adminCategories := protected.Group("/categories", middleware.AdminMiddleware())
		{
			adminCategories.POST("", categoryH.CreateCategory)
			adminCategories.PUT("/:categoryId", categoryH.UpdateCategory)
			adminCategories.DELETE("/:categoryId", categoryH.DeleteCategory)
		}

		adminNominees := protected.Group("/nominees", middleware.AdminMiddleware())
		{
			adminNominees.POST("", nomineeH.CreateNominee)
			adminNominees.PUT("/:id", nomineeH.UpdateNominee)
			adminNominees.DELETE("/:id", nomineeH.DeleteNominee)
		}

		/*	nomineeCategories := protected.Group("/nominees/:nominee_id/categories")
			{
				nomineeCategories.POST("", nomineeCategoryH.AddCategory)
				//nomineeCategories.DELETE("/:category_id", nomineeCategoryH.RemoveCategory)
				nomineeCategories.PUT("", nomineeCategoryH.SetCategories)
				nomineeCategories.GET("", nomineeCategoryH.GetCategories)
			}
		*/
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
