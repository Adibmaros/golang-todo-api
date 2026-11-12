package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Adibmaros/todo-api/internal/model"
)

type TodoHandler struct {
	mu     sync.Mutex
	todos  map[int64]model.Todo
	nextID int64
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todos:  make(map[int64]model.Todo),
		nextID: 1,
	}
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	todos := make([]model.Todo, 0, len(h.todos))

	for _, todo := range h.todos {
		todos = append(todos, todo)
	}

	JSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		Error(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	todo, exists := h.todos[id]

	if !exists {
		Error(w, http.StatusNotFound, "todo not found")
		return
	}

	JSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, " gagal tambah data todo")
		return
	}

	if input.Title == "" {
		Error(w, http.StatusBadRequest, "title must be input")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	todo := model.Todo{
		ID:        h.nextID,
		Title:     input.Title,
		Done:      false,
		CreatedAt: time.Now(),
	}

	h.todos[todo.ID] = todo
	h.nextID++

	JSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	var input struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.Title == "" {
		Error(w, http.StatusBadRequest, "title is required")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	todo, exists := h.todos[id]
	if !exists {
		Error(w, http.StatusNotFound, "todo not found")
		return
	}

	todo.Title = input.Title
	todo.Done = input.Done

	h.todos[id] = todo

	JSON(w, http.StatusOK, todo)
}

// DELETE /todos/{id}
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.todos[id]; !exists {
		Error(w, http.StatusNotFound, "todo not found")
		return
	}

	delete(h.todos, id)

	w.WriteHeader(http.StatusNoContent)
}
