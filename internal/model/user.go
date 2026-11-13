package model

import "time"

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Todos     []Todo    `json:"todos,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
}
