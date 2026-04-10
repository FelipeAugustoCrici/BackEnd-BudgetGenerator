package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Company   string         `json:"company"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Notes     string         `json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Client) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Contract struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ClientID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"client_id"`
	Client      *Client        `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	BudgetID    *uuid.UUID     `gorm:"type:uuid" json:"budget_id,omitempty"`
	Value       float64        `gorm:"type:decimal(12,2)" json:"value"`
	Status      string         `gorm:"default:draft" json:"status"`
	StartDate   *time.Time     `json:"start_date,omitempty"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	AutoRenew   bool           `gorm:"default:false" json:"auto_renew"`
	Description string         `json:"description"`
	SentAt      *time.Time     `json:"sent_at,omitempty"`
	ViewedAt    *time.Time     `json:"viewed_at,omitempty"`
	SignedAt    *time.Time     `json:"signed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Contract) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type ContractEvent struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ContractID uuid.UUID  `gorm:"type:uuid;not null;index" json:"contract_id"`
	Type       string     `gorm:"not null" json:"type"`
	Metadata   string     `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (e *ContractEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
