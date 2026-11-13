package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Adibmaros/todo-api/internal/service"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: s,
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

	todo, err := h.service.Create(input.Title)
	if errors.Is(err, service.ErrTitleRequired) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal membuat data",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

	todo, err := h.service.GetByID(id)
	if errors.Is(err, service.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
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

	todo, err := h.service.Update(id, input.Title, input.Done)
	if errors.Is(err, service.ErrTitleRequired) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, service.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal menyimpan data",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	err := h.service.Delete(id)
	if errors.Is(err, service.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal menghapus data",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
