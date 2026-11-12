# Golang Todo API

Proyek belajar membangun REST API dengan Go dari nol menggunakan standard library `net/http` (`http.ServeMux`), MySQL, dan GORM.

## Tech Stack

- **Go** (`1.25.5`)
- **HTTP Router:** `net/http` (`http.ServeMux`)
- **ORM:** [GORM](https://gorm.io/) (`gorm.io/gorm`)
- **Database:** MySQL (`gorm.io/driver/mysql`)

## Struktur Proyek

```text
todo-api/
├── cmd/
│   └── main.go                 # Entry point aplikasi & routing
├── internal/
│   ├── database/
│   │   └── mysql.go            # Koneksi database MySQL via GORM
│   ├── handler/
│   │   ├── health.go           # Handler health check
│   │   └── todo.go             # Handler CRUD Todo
│   ├── model/
│   │   └── todo.go             # Struct & schema model Todo
│   └── response/
│       └── response.go         # Helper response JSON & error
├── migrations/
│   └── init.sql                # SQL pembuatan database awal
├── .gitignore
├── go.mod
└── go.sum
```

## Cara Menjalankan

1. **Siapkan database MySQL** (pastikan MySQL berjalan di `127.0.0.1:3306`), lalu jalankan query di `migrations/init.sql`:
   ```sql
   CREATE DATABASE todo_db;
   ```
2. **Sesuaikan konfigurasi koneksi** di `internal/database/mysql.go` jika username, password, atau port MySQL berbeda (default saat ini menggunakan user `root` tanpa password).
3. **Jalankan server:**
   ```bash
   go run ./cmd/main.go
   ```
   Server akan berjalan di `http://localhost:8080` dan GORM `AutoMigrate` akan otomatis membuat tabel `todos`.

## Daftar Endpoint

| Method | Endpoint | Deskripsi | Contoh Body |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Cek status server | - |
| `GET` | `/todos` | Ambil semua todo | - |
| `GET` | `/todos/{id}` | Ambil detail todo berdasarkan ID | - |
| `POST` | `/todos` | Tambah todo baru | `{"title": "Belajar Go"}` |
| `PUT` | `/todos/{id}` | Update todo berdasarkan ID | `{"title": "Belajar Go", "done": true}` |
| `DELETE` | `/todos/{id}` | Hapus todo berdasarkan ID | - |
