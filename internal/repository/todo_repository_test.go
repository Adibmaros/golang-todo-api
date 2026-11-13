package repository

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Adibmaros/todo-api/internal/dto"
	"github.com/Adibmaros/todo-api/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dbName := os.Getenv("TEST_DB_NAME")
		if dbName == "" {
			t.Skip("lewati integration test: set TEST_DB_NAME atau TEST_DB_DSN untuk menjalankan test database MySQL")
		}

		host := envOrDefault("TEST_DB_HOST", "127.0.0.1")
		port := envOrDefault("TEST_DB_PORT", "3306")
		user := envOrDefault("TEST_DB_USER", "root")
		pass := os.Getenv("TEST_DB_PASS")

		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, pass, host, port, dbName,
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("gagal terhubung ke database test MySQL: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		t.Fatalf("gagal migrasi tabel test: %v", err)
	}

	if err := db.Exec("DELETE FROM todos").Error; err != nil {
		t.Fatalf("gagal membersihkan tabel todos: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Exec("DELETE FROM todos").Error
	})

	return db
}

func envOrDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func TestTodoRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)

	todo1 := &model.Todo{UserID: 1, Title: "Belajar A", Done: false}
	todo2 := &model.Todo{UserID: 1, Title: "Belajar B", Done: true}
	otherUserTodo := &model.Todo{UserID: 2, Title: "Milik User 2", Done: false}

	for _, td := range []*model.Todo{todo1, todo2, otherUserTodo} {
		if err := repo.Create(td); err != nil {
			t.Fatalf("gagal create todo: %v", err)
		}
	}

	// FindAllByUserID dengan filter done=true
	doneTrue := true
	todos, total, err := repo.FindAllByUserID(1, dto.TodoQuery{
		Page:  1,
		Limit: 10,
		Done:  &doneTrue,
		Sort:  "title",
	})
	if err != nil {
		t.Fatalf("gagal FindAllByUserID: %v", err)
	}
	if total != 1 || len(todos) != 1 || todos[0].ID != todo2.ID {
		t.Fatalf("hasil filter tidak sesuai: total=%d, todos=%+v", total, todos)
	}

	// Isolasi user: User 1 tidak bisa mengakses todo milik User 2
	_, err = repo.FindByIDAndUserID(otherUserTodo.ID, 1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("mengharapkan ErrNotFound untuk todo milik user lain, dapat: %v", err)
	}

	// Update
	todo1.Title = "Belajar A Updated"
	todo1.Done = true
	if err := repo.Update(todo1); err != nil {
		t.Fatalf("gagal update todo: %v", err)
	}

	updated, err := repo.FindByIDAndUserID(todo1.ID, 1)
	if err != nil {
		t.Fatalf("gagal mengambil todo setelah update: %v", err)
	}
	if updated.Title != "Belajar A Updated" || !updated.Done {
		t.Fatalf("data setelah update tidak sesuai: %+v", updated)
	}

	// Delete
	if err := repo.DeleteByIDAndUserID(todo1.ID, 1); err != nil {
		t.Fatalf("gagal delete todo: %v", err)
	}
	if err := repo.DeleteByIDAndUserID(todo1.ID, 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("mengharapkan ErrNotFound saat hapus ulang, dapat: %v", err)
	}
}
