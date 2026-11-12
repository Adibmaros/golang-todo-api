package model

import "time"

type Todo struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"not null"`
	Done      bool      `json:"done" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
