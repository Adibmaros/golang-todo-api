package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/response"

	"gorm.io/gorm"
)

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{
		db: db,
	}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// deskripsiin input type

	var input struct {
		Title string `json:"title"`
	}

	// ambil request dari body

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "input tidak valid")
		return
	}

	if input.Title == "" {
		response.Error(w, http.StatusBadRequest, "title harus diisi")
		return
	}

	todo := model.Todo{
		Title: input.Title,
		Done:  false,
	}

	if err := h.db.Create(&todo).Error; err != nil {
		response.Error(w, http.StatusBadRequest, "gagal membuat todo")
		return
	}

	response.JSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	var todos []model.Todo

	if err := h.db.Find(&todos).Error; err != nil {
		response.Error(w, http.StatusBadRequest, "gagal ambil data dari db")
		return
	}

	response.JSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	var todo model.Todo

	err = h.db.First(&todo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Error(w, http.StatusNotFound, "data tidak ditemukan")
		return
	}

	if err != nil {
		response.Error(w, http.StatusBadRequest, "gagal ambil data todo")
		return
	}

	response.JSON(w, http.StatusOK, todo)

}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id type")
		return
	}

	var input struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "request body gak valid")
		return
	}

	var todo model.Todo

	if err := h.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal mengambil todo")
		return
	}

	result := h.db.Model(&todo).Updates(map[string]interface{}{
		"title": input.Title,
		"done":  input.Done,
	})

	if result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "gagal simpan data")
		return
	}

	todo.Title = input.Title
	todo.Done = input.Done

	response.JSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	result := h.db.Delete(&model.Todo{}, id)

	if result.Error != nil {
		response.Error(w, http.StatusBadRequest, "gagal menghapus data")
		return
	}

	if result.RowsAffected == 0 {
		response.Error(w, http.StatusNotFound, "data tidak ditemukan")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
