package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Adibmaros/todo-api/internal/dto"
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
	var input dto.CreateTodoRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	todo, err := h.service.Create(userID, input.Title)
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

func parseTodoQuery(c *gin.Context) (dto.TodoQuery, bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page harus berupa angka >= 1"})
		return dto.TodoQuery{}, false
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit harus berupa angka >= 1"})
		return dto.TodoQuery{}, false
	}
	if limit > 100 {
		limit = 100
	}

	var donePtr *bool
	if doneStr := c.Query("done"); doneStr != "" {
		doneVal, err := strconv.ParseBool(doneStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "done harus bernilai true atau false"})
			return dto.TodoQuery{}, false
		}
		donePtr = &doneVal
	}

	sort := c.DefaultQuery("sort", "created_at")

	return dto.TodoQuery{
		Page:  page,
		Limit: limit,
		Done:  donePtr,
		Sort:  sort,
	}, true
}

func (h *TodoHandler) List(c *gin.Context) {
	query, ok := parseTodoQuery(c)
	if !ok {
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.service.List(userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, res)
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

	userID := c.GetInt64("user_id")
	todo, err := h.service.GetByID(id, userID)
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

	var input dto.UpdateTodoRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	todo, err := h.service.Update(id, userID, input.Title, input.Done)
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

	userID := c.GetInt64("user_id")
	err := h.service.Delete(id, userID)
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
