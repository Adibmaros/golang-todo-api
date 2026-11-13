package dto

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=3,max=100"`
}

type UpdateTodoRequest struct {
	Title string `json:"title" binding:"required,min=3,max=10"`
	Done  bool   `json:"done"`
}
