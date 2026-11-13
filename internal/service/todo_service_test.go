package service

import (
	"errors"
	"testing"

	"github.com/Adibmaros/todo-api/internal/dto"
	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/repository"
)

type mockTodoRepository struct {
	createFn            func(todo *model.Todo) error
	findAllByUserIDFn   func(userID int64, query dto.TodoQuery) ([]model.Todo, int64, error)
	findByIDAndUserIDFn func(id, userID int64) (*model.Todo, error)
	updateFn            func(todo *model.Todo) error
	deleteByIDAndUserFn func(id, userID int64) error
}

func (m *mockTodoRepository) Create(todo *model.Todo) error {
	if m.createFn != nil {
		return m.createFn(todo)
	}
	return nil
}

func (m *mockTodoRepository) FindAllByUserID(userID int64, query dto.TodoQuery) ([]model.Todo, int64, error) {
	if m.findAllByUserIDFn != nil {
		return m.findAllByUserIDFn(userID, query)
	}
	return nil, 0, nil
}

func (m *mockTodoRepository) FindByIDAndUserID(id, userID int64) (*model.Todo, error) {
	if m.findByIDAndUserIDFn != nil {
		return m.findByIDAndUserIDFn(id, userID)
	}
	return nil, nil
}

func (m *mockTodoRepository) Update(todo *model.Todo) error {
	if m.updateFn != nil {
		return m.updateFn(todo)
	}
	return nil
}

func (m *mockTodoRepository) DeleteByIDAndUserID(id, userID int64) error {
	if m.deleteByIDAndUserFn != nil {
		return m.deleteByIDAndUserFn(id, userID)
	}
	return nil
}

func TestTodoService_Create(t *testing.T) {
	t.Run("berhasil membuat todo", func(t *testing.T) {
		repo := &mockTodoRepository{
			createFn: func(todo *model.Todo) error {
				todo.ID = 10
				return nil
			},
		}
		svc := NewTodoService(repo)

		todo, err := svc.Create(7, "Belajar Go Testing")
		if err != nil {
			t.Fatalf("tidak mengharapkan error, dapat: %v", err)
		}
		if todo.ID != 10 || todo.UserID != 7 || todo.Title != "Belajar Go Testing" || todo.Done {
			t.Fatalf("data todo tidak sesuai: %+v", todo)
		}
	})

	t.Run("gagal jika title kosong atau hanya spasi", func(t *testing.T) {
		svc := NewTodoService(&mockTodoRepository{})

		_, err := svc.Create(7, "   ")
		if !errors.Is(err, ErrTitleRequired) {
			t.Fatalf("mengharapkan ErrTitleRequired, dapat: %v", err)
		}
	})
}

func TestTodoService_List(t *testing.T) {
	t.Run("menghitung total_pages dengan benar", func(t *testing.T) {
		repo := &mockTodoRepository{
			findAllByUserIDFn: func(userID int64, query dto.TodoQuery) ([]model.Todo, int64, error) {
				if userID != 5 {
					t.Fatalf("mengharapkan userID 5, dapat: %d", userID)
				}
				return []model.Todo{{ID: 1, UserID: 5, Title: "Satu"}}, 21, nil
			},
		}
		svc := NewTodoService(repo)

		res, err := svc.List(5, dto.TodoQuery{Page: 1, Limit: 10})
		if err != nil {
			t.Fatalf("tidak mengharapkan error, dapat: %v", err)
		}
		if res.Meta.TotalData != 21 || res.Meta.TotalPage != 3 {
			t.Fatalf("meta pagination tidak sesuai: %+v", res.Meta)
		}
		if len(res.Data) != 1 {
			t.Fatalf("mengharapkan 1 data, dapat: %d", len(res.Data))
		}
	})

	t.Run("total_pages nol saat data kosong", func(t *testing.T) {
		repo := &mockTodoRepository{
			findAllByUserIDFn: func(userID int64, query dto.TodoQuery) ([]model.Todo, int64, error) {
				return []model.Todo{}, 0, nil
			},
		}
		svc := NewTodoService(repo)

		res, err := svc.List(5, dto.TodoQuery{Page: 1, Limit: 10})
		if err != nil {
			t.Fatalf("tidak mengharapkan error, dapat: %v", err)
		}
		if res.Meta.TotalPage != 0 {
			t.Fatalf("mengharapkan TotalPage 0, dapat: %d", res.Meta.TotalPage)
		}
	})
}

func TestTodoService_GetByID(t *testing.T) {
	t.Run("berhasil mengambil todo", func(t *testing.T) {
		repo := &mockTodoRepository{
			findByIDAndUserIDFn: func(id, userID int64) (*model.Todo, error) {
				return &model.Todo{ID: id, UserID: userID, Title: "Test"}, nil
			},
		}
		svc := NewTodoService(repo)

		todo, err := svc.GetByID(3, 1)
		if err != nil {
			t.Fatalf("tidak mengharapkan error, dapat: %v", err)
		}
		if todo.ID != 3 || todo.UserID != 1 {
			t.Fatalf("todo tidak sesuai: %+v", todo)
		}
	})

	t.Run("mengembalikan ErrTodoNotFound saat tidak ditemukan", func(t *testing.T) {
		repo := &mockTodoRepository{
			findByIDAndUserIDFn: func(id, userID int64) (*model.Todo, error) {
				return nil, repository.ErrNotFound
			},
		}
		svc := NewTodoService(repo)

		_, err := svc.GetByID(99, 1)
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("mengharapkan ErrTodoNotFound, dapat: %v", err)
		}
	})
}

func TestTodoService_Update(t *testing.T) {
	t.Run("berhasil mengupdate todo", func(t *testing.T) {
		repo := &mockTodoRepository{
			findByIDAndUserIDFn: func(id, userID int64) (*model.Todo, error) {
				return &model.Todo{ID: id, UserID: userID, Title: "Lama", Done: false}, nil
			},
			updateFn: func(todo *model.Todo) error {
				return nil
			},
		}
		svc := NewTodoService(repo)

		todo, err := svc.Update(2, 1, "Baru", true)
		if err != nil {
			t.Fatalf("tidak mengharapkan error, dapat: %v", err)
		}
		if todo.Title != "Baru" || !todo.Done {
			t.Fatalf("todo tidak terupdate: %+v", todo)
		}
	})

	t.Run("gagal jika title kosong", func(t *testing.T) {
		svc := NewTodoService(&mockTodoRepository{})

		_, err := svc.Update(2, 1, "", true)
		if !errors.Is(err, ErrTitleRequired) {
			t.Fatalf("mengharapkan ErrTitleRequired, dapat: %v", err)
		}
	})

	t.Run("mengembalikan ErrTodoNotFound jika tidak ditemukan", func(t *testing.T) {
		repo := &mockTodoRepository{
			findByIDAndUserIDFn: func(id, userID int64) (*model.Todo, error) {
				return nil, repository.ErrNotFound
			},
		}
		svc := NewTodoService(repo)

		_, err := svc.Update(99, 1, "Baru", false)
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("mengharapkan ErrTodoNotFound, dapat: %v", err)
		}
	})
}

func TestTodoService_Delete(t *testing.T) {
	t.Run("berhasil menghapus todo", func(t *testing.T) {
		repo := &mockTodoRepository{
			deleteByIDAndUserFn: func(id, userID int64) error {
				return nil
			},
		}
		svc := NewTodoService(repo)

		if err := svc.Delete(1, 1); err != nil {
			t.Fatalf("tidak mengharapkan error, dapat: %v", err)
		}
	})

	t.Run("mengembalikan ErrTodoNotFound saat tidak ditemukan", func(t *testing.T) {
		repo := &mockTodoRepository{
			deleteByIDAndUserFn: func(id, userID int64) error {
				return repository.ErrNotFound
			},
		}
		svc := NewTodoService(repo)

		err := svc.Delete(99, 1)
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("mengharapkan ErrTodoNotFound, dapat: %v", err)
		}
	})
}

