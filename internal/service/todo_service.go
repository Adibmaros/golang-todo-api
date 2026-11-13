package service

import (
	"errors"
	"strings"

	"github.com/Adibmaros/todo-api/internal/dto"
	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/repository"
)

var (
	ErrTitleRequired = errors.New("title harus diisi")
	ErrTodoNotFound  = errors.New("data tidak ditemukan")
)

type TodoService interface {
	Create(userID int64, title string) (*model.Todo, error)
	List(userID int64, query dto.TodoQuery) (*dto.TodoListResponse, error)
	GetByID(id, userID int64) (*model.Todo, error)
	Update(id, userID int64, title string, done bool) (*model.Todo, error)
	Delete(id, userID int64) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(userID int64, title string) (*model.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrTitleRequired
	}

	todo := &model.Todo{
		UserID: userID,
		Title:  title,
		Done:   false,
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) List(userID int64, query dto.TodoQuery) (*dto.TodoListResponse, error) {
	todos, totalData, err := s.repo.FindAllByUserID(userID, query)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if totalData > 0 {
		totalPages = int((totalData + int64(query.Limit) - 1) / int64(query.Limit))
	}

	return &dto.TodoListResponse{
		Data: todos,
		Meta: dto.PaginationMeta{
			Page:      query.Page,
			Limit:     query.Limit,
			TotalData: totalData,
			TotalPage: totalPages,
		},
	}, nil
}

func (s *todoService) GetByID(id, userID int64) (*model.Todo, error) {
	todo, err := s.repo.FindByIDAndUserID(id, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoService) Update(id, userID int64, title string, done bool) (*model.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrTitleRequired
	}

	todo, err := s.repo.FindByIDAndUserID(id, userID)
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

func (s *todoService) Delete(id, userID int64) error {
	err := s.repo.DeleteByIDAndUserID(id, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrTodoNotFound
	}
	return err
}
