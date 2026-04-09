package model

import "time"

type Settings struct {
	ID              string    `gorm:"primaryKey;type:uuid" json:"id"`
	UserID          string    `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Name            string    `json:"name"`
	Logo            string    `gorm:"type:text" json:"logo"`
	Instagram       string    `json:"instagram"`
	Facebook        string    `json:"facebook"`
	Website         string    `json:"website"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	FontFamily      string    `json:"fontFamily"`
	FontURL         string    `json:"fontUrl"`
	ShowNameOnHeader bool     `gorm:"default:true" json:"showNameOnHeader"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
