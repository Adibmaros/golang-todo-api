package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Adibmaros/todo-api/docs"
	"github.com/Adibmaros/todo-api/internal/database"
	"github.com/Adibmaros/todo-api/internal/handler"
	"github.com/Adibmaros/todo-api/internal/middleware"
	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/repository"
	"github.com/Adibmaros/todo-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title                       Golang Todo API
// @version                     1.0
// @description                 REST API Todo dengan Gin, MySQL, GORM, dan autentikasi JWT.
// @BasePath                    /
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Masukkan token dengan format: Bearer <token>

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Info("file .env tidak ditemukan, menggunakan environment variable sistem")
	}

	db := database.Connect()

	if err := db.AutoMigrate(&model.Todo{}, &model.User{}); err != nil {
		slog.Error("migrasi gagal", "error", err)
		os.Exit(1)
	}

	healthHandler := handler.NewHealthHandler(db)
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(60, time.Minute))

	r.GET("/health", healthHandler.Check)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	todos := r.Group("/todos", middleware.Auth())
	{
		todos.GET("", todoHandler.List)
		todos.GET("/:id", todoHandler.Get)
		todos.POST("", todoHandler.Create)
		todos.PUT("/:id", todoHandler.Update)
		todos.DELETE("/:id", todoHandler.Delete)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server berjalan", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server gagal berjalan", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("menerima sinyal shutdown, menghentikan server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown server gagal", "error", err)
		os.Exit(1)
	}

	slog.Info("server berhenti dengan aman")
}
