package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"code-challenge/internal/auth"
	"code-challenge/internal/config"
	"code-challenge/internal/form"
	"code-challenge/internal/user"
	"code-challenge/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Version = "1.0.0"

var imageDir = "./tmp"

func main() {
	// Create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// Load application configurations
	configFile := "local"
	env := os.Getenv("APP_ENVIRONMENT")
	if env != "" {
		configFile = env
	}
	logger.Infof("Loading config from %s", configFile)

	appConfig, err := config.Load(configFile)
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  appConfig.DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "codechallenge.",
			SingularTable: false,
		}})
	if err != nil {
		logger.Fatalf("Failed to initialize database connection: %v", err)
	}

	// Initialize routes
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo, logger)
	userHandler := user.NewHandler(userService, logger)

	jwtService := auth.NewJWTService(appConfig.JWTSigningKey, appConfig.JWTExpiration)
	authService := auth.NewService(userRepo, jwtService, logger)
	authHandler := auth.NewHandler(authService, logger)

	authMiddleware := auth.Middleware(jwtService, logger)

	userRouterGroup := router.Group("/api/v0/users")
	user.RegisterHandlers(userRouterGroup, authMiddleware, userHandler)

	authRouterGroup := router.Group("/api/v0/auth")
	router.POST("/api/v0/auth")
	auth.RegisterHandlers(authRouterGroup, authHandler)

	formRepo := form.NewRepo(db)
	formService := form.NewService(imageDir, formRepo, logger)
	formHandler := form.NewHandler(formService, logger)
	formRouterGroup := router.Group("")
	form.RegisterHandlers(formRouterGroup, authMiddleware, formHandler)

	// HTTP Server
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(appConfig.ServerPort),
		Handler: router,
	}

	// Initializing the server in a goroutine so that it won't block the gracefully shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	logger.Info("Server is running...")

	<-quit
	logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}
	logger.Info("Server exiting")
}
