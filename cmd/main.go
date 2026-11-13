package main

import (
	"log"

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

	if err := godotenv.Load(); err != nil {
		log.Println("file .env tidak ditemukan")
	}

	db := database.Connect()

	if err := db.AutoMigrate(&model.Todo{}, &model.User{}); err != nil {
		log.Fatal("migrasi gagal", err)
	}

	healthHandler := handler.NewHealthHandler()
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

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

	log.Println("berjalan di port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("server gagal berjalan", err)
	}
}
