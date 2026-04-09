package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type TemplateBlocks []TemplateBlock

type TemplateBlock struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Align   string `json:"align,omitempty"`
	Visible bool   `json:"visible"`
}

func (t TemplateBlocks) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *TemplateBlocks) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, t)
}

type Template struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"userId"`
	Name      string         `gorm:"not null" json:"name"`
	IsDefault bool           `gorm:"default:false" json:"isDefault"`
	Blocks    TemplateBlocks `gorm:"type:jsonb" json:"blocks"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
