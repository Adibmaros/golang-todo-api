package repository

import (
	"errors"

	"github.com/Adibmaros/todo-api/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("data tidak ditemukan")

type TodoRepository interface {
	Create(todo *model.Todo) error
	FindAll() ([]model.Todo, error)
	FindByID(id int64) (*model.Todo, error)
	Update(todo *model.Todo) error
	Delete(id int64) error
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

func (r *todoRepository) FindAll() ([]model.Todo, error) {
	var todos []model.Todo
	if err := r.db.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}
func (r *todoRepository) FindByID(id int64) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.First(&todo, id).Error
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
func (r *todoRepository) Delete(id int64) error {
	result := r.db.Delete(&model.Todo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
