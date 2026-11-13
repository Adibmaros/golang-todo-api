package repository

import (
	"errors"

	"github.com/Adibmaros/todo-api/internal/dto"
	"github.com/Adibmaros/todo-api/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("data tidak ditemukan")

var allowedSortColumns = map[string]string{
	"created_at":  "created_at DESC",
	"created_asc": "created_at ASC",
	"title":       "title ASC",
	"done":        "done ASC",
}

type TodoRepository interface {
	Create(todo *model.Todo) error
	FindAllByUserID(userID int64, query dto.TodoQuery) ([]model.Todo, int64, error)
	FindByIDAndUserID(id, userID int64) (*model.Todo, error)
	Update(todo *model.Todo) error
	DeleteByIDAndUserID(id, userID int64) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) FindAllByUserID(userID int64, query dto.TodoQuery) ([]model.Todo, int64, error) {
	db := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)

	if query.Done != nil {
		db = db.Where("done = ?", *query.Done)
	}

	var totalData int64
	if err := db.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	orderBy, ok := allowedSortColumns[query.Sort]
	if !ok {
		orderBy = "created_at DESC"
	}

	offset := (query.Page - 1) * query.Limit

	var todos []model.Todo
	if err := db.Order(orderBy).Limit(query.Limit).Offset(offset).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, totalData, nil
}
func (r *todoRepository) FindByIDAndUserID(id, userID int64) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}
func (r *todoRepository) Update(todo *model.Todo) error {
	// Pakai map agar nilai zero-value seperti done=false tetap ter-update oleh GORM
	return r.db.Model(todo).Updates(map[string]any{
		"title": todo.Title,
		"done":  todo.Done,
	}).Error
}
func (r *todoRepository) DeleteByIDAndUserID(id, userID int64) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
