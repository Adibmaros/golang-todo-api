package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Adibmaros/todo-api/internal/dto"
	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/service"
	"github.com/gin-gonic/gin"
)

type mockTodoService struct {
	createFn  func(userID int64, title string) (*model.Todo, error)
	listFn    func(userID int64, query dto.TodoQuery) (*dto.TodoListResponse, error)
	getByIDFn func(id, userID int64) (*model.Todo, error)
	updateFn  func(id, userID int64, title string, done bool) (*model.Todo, error)
	deleteFn  func(id, userID int64) error
}

func (m *mockTodoService) Create(userID int64, title string) (*model.Todo, error) {
	if m.createFn != nil {
		return m.createFn(userID, title)
	}
	return nil, nil
}

func (m *mockTodoService) List(userID int64, query dto.TodoQuery) (*dto.TodoListResponse, error) {
	if m.listFn != nil {
		return m.listFn(userID, query)
	}
	return nil, nil
}

func (m *mockTodoService) GetByID(id, userID int64) (*model.Todo, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id, userID)
	}
	return nil, nil
}

func (m *mockTodoService) Update(id, userID int64, title string, done bool) (*model.Todo, error) {
	if m.updateFn != nil {
		return m.updateFn(id, userID, title, done)
	}
	return nil, nil
}

func (m *mockTodoService) Delete(id, userID int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id, userID)
	}
	return nil
}

func setupTodoTestRouter(svc service.TodoService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewTodoHandler(svc)

	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(42))
		c.Next()
	})

	r.GET("/todos", h.List)
	r.GET("/todos/:id", h.Get)
	r.POST("/todos", h.Create)
	r.PUT("/todos/:id", h.Update)
	r.DELETE("/todos/:id", h.Delete)

	return r
}

func TestTodoHandler_Create(t *testing.T) {
	t.Run("201 created saat request valid", func(t *testing.T) {
		svc := &mockTodoService{
			createFn: func(userID int64, title string) (*model.Todo, error) {
				if userID != 42 {
					t.Fatalf("mengharapkan userID 42, dapat: %d", userID)
				}
				return &model.Todo{ID: 1, UserID: userID, Title: title}, nil
			},
		}
		r := setupTodoTestRouter(svc)

		body := bytes.NewBufferString(`{"title":"Belajar Handler Test"}`)
		req := httptest.NewRequest(http.MethodPost, "/todos", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("mengharapkan status 201, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		var got model.Todo
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("gagal decode response: %v", err)
		}
		if got.Title != "Belajar Handler Test" {
			t.Fatalf("title tidak sesuai: %s", got.Title)
		}
	})

	t.Run("400 bad request saat validasi gagal", func(t *testing.T) {
		r := setupTodoTestRouter(&mockTodoService{})

		body := bytes.NewBufferString(`{"title":"ab"}`)
		req := httptest.NewRequest(http.MethodPost, "/todos", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("mengharapkan status 400, dapat: %d", w.Code)
		}
	})
}

func TestTodoHandler_List(t *testing.T) {
	t.Run("200 ok dengan query pagination dan filter", func(t *testing.T) {
		svc := &mockTodoService{
			listFn: func(userID int64, query dto.TodoQuery) (*dto.TodoListResponse, error) {
				if query.Page != 2 || query.Limit != 5 || query.Done == nil || !*query.Done || query.Sort != "title" {
					t.Fatalf("query tidak sesuai: %+v", query)
				}
				return &dto.TodoListResponse{
					Data: []model.Todo{{ID: 1, UserID: userID, Title: "A", Done: true}},
					Meta: dto.PaginationMeta{Page: 2, Limit: 5, TotalData: 6, TotalPage: 2},
				}, nil
			},
		}
		r := setupTodoTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos?page=2&limit=5&done=true&sort=title", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("mengharapkan status 200, dapat: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("400 bad request jika query done tidak valid", func(t *testing.T) {
		r := setupTodoTestRouter(&mockTodoService{})

		req := httptest.NewRequest(http.MethodGet, "/todos?done=bukanbool", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("mengharapkan status 400, dapat: %d", w.Code)
		}
	})
}

func TestTodoHandler_Get(t *testing.T) {
	t.Run("404 not found saat todo tidak ada", func(t *testing.T) {
		svc := &mockTodoService{
			getByIDFn: func(id, userID int64) (*model.Todo, error) {
				return nil, service.ErrTodoNotFound
			},
		}
		r := setupTodoTestRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/todos/99", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("mengharapkan status 404, dapat: %d", w.Code)
		}
	})

	t.Run("400 bad request jika id bukan angka", func(t *testing.T) {
		r := setupTodoTestRouter(&mockTodoService{})

		req := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("mengharapkan status 400, dapat: %d", w.Code)
		}
	})
}

func TestTodoHandler_Update(t *testing.T) {
	t.Run("200 ok saat update berhasil", func(t *testing.T) {
		svc := &mockTodoService{
			updateFn: func(id, userID int64, title string, done bool) (*model.Todo, error) {
				return &model.Todo{ID: id, UserID: userID, Title: title, Done: done}, nil
			},
		}
		r := setupTodoTestRouter(svc)

		body := bytes.NewBufferString(`{"title":"Update","done":true}`)
		req := httptest.NewRequest(http.MethodPut, "/todos/1", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("mengharapkan status 200, dapat: %d, body: %s", w.Code, w.Body.String())
		}
	})
}

func TestTodoHandler_Delete(t *testing.T) {
	t.Run("204 no content saat delete berhasil", func(t *testing.T) {
		svc := &mockTodoService{
			deleteFn: func(id, userID int64) error {
				return nil
			},
		}
		r := setupTodoTestRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("mengharapkan status 204, dapat: %d", w.Code)
		}
	})
}

