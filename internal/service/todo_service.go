package service

import (
	"errors"
	"strings"

	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/repository"
)

var (
	ErrTitleRequired = errors.New("title harus diisi")
	ErrTodoNotFound  = errors.New("data tidak ditemukan")
)

type TodoService interface {
	Create(title string) (*model.Todo, error)
	List() ([]model.Todo, error)
	GetByID(id int64) (*model.Todo, error)
	Update(id int64, title string, done bool) (*model.Todo, error)
	Delete(id int64) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(title string) (*model.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrTitleRequired
	}

	todo := &model.Todo{
		Title: title,
		Done:  false,
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) List() ([]model.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) GetByID(id int64) (*model.Todo, error) {
	todo, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoService) Update(id int64, title string, done bool) (*model.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrTitleRequired
	}

	todo, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}

	todo.Title = title
	todo.Done = done

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrTodoNotFound
	}
	return err
}
