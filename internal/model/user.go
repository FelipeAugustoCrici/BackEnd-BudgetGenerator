package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
