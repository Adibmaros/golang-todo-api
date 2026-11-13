package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/gin-gonic/gin"
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

func (h *TodoHandler) Create(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body tidak valid",
		})
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title harus diisi",
		})
		return
	}

	todo := model.Todo{
		Title: input.Title,
		Done:  false,
	}

	if err := h.db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gagal membuat data",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) List(c *gin.Context) {
	var todos []model.Todo

	if err := h.db.Find(&todos).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func parseIDParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id tidak valid",
		})
		return 0, false
	}
	return id, true
}

func (h *TodoHandler) Get(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var todo model.Todo
	err := h.db.First(&todo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "data tidak ditemukan",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var input struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body tidak valid",
		})
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title harus diisi",
		})
		return
	}

	var todo model.Todo
	if err := h.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "data tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil data",
		})
		return
	}

	result := h.db.Model(&todo).Updates(map[string]any{
		"title": input.Title,
		"done":  input.Done,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal menyimpan data",
		})
		return
	}

	todo.Title = input.Title
	todo.Done = input.Done

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	result := h.db.Delete(&model.Todo{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal menghapus data",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "data tidak ditemukan",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
