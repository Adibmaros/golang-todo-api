package main

import (
	"log"

	"github.com/Adibmaros/todo-api/internal/database"
	"github.com/Adibmaros/todo-api/internal/handler"
	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("file .env tidak ditemukan")
	}

	db := database.Connect()

	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		log.Fatal("migrasi gagal", err)
	}

	healthHandler := handler.NewHealthHandler()
	todoHandler := handler.NewTodoHandler(db)

	r := gin.Default()

	r.GET("/health", healthHandler.Check)
	r.GET("/todos", todoHandler.List)
	r.GET("/todos/:id", todoHandler.Get)
	r.POST("/todos", todoHandler.Create)
	r.PUT("/todos/:id", todoHandler.Update)
	r.DELETE("/todos/:id", todoHandler.Delete)

	log.Println("berjalan di port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("server gagal berjalan", err)
	}
}
