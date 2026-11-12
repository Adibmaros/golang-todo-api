package main

import (
	"log"
	"net/http"

	"github.com/Adibmaros/todo-api/internal/handler"
)

func main() {

	healthHandler := handler.NewHealthHandler()
	todoHandler := handler.NewTodoHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler.Check)
	mux.HandleFunc("GET /todos", todoHandler.List)
	mux.HandleFunc("GET /todos/{id}", todoHandler.Get)
	mux.HandleFunc("POST /todos", todoHandler.Create)
	mux.HandleFunc("PUT /todos/{id}", todoHandler.Update)
	mux.HandleFunc("DELETE /todos/{id}", todoHandler.Delete)

	log.Println("berjalan di port 8080")

	http.ListenAndServe(":8080", mux)

}
