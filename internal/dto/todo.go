package dto

import "github.com/Adibmaros/todo-api/internal/model"

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=3,max=100"`
}

type UpdateTodoRequest struct {
	Title string `json:"title" binding:"required,min=3,max=10"`
	Done  bool   `json:"done"`
}

type TodoQuery struct {
	Page  int
	Limit int
	Done  *bool
	Sort  string
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_pages"`
}

type TodoListResponse struct {
	Data []model.Todo   `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
