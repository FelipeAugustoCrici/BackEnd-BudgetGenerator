package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// QuoteItems is a JSON column
type QuoteItems []QuoteItem

type QuoteItem struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unitPrice"`
	EstimateHours float64 `json:"estimateHours,omitempty"`
	ItemStatus    string  `json:"itemStatus,omitempty"`
}

func (q QuoteItems) Value() (driver.Value, error) {
	return json.Marshal(q)
}

func (q *QuoteItems) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, q)
}

type Quote struct {
	ID           string     `gorm:"primaryKey;type:uuid" json:"id"`
	UserID       string     `gorm:"type:uuid;not null;index" json:"userId"`
	ClientName   string     `json:"clientName"`
	Date         string     `json:"date"`
	Notes        string     `json:"notes"`
	Status       string     `gorm:"default:draft" json:"status"`
	Items        QuoteItems `gorm:"type:jsonb" json:"items"`
	Discount     float64    `json:"discount"`
	DiscountType string     `gorm:"default:percent" json:"discountType"`
	TemplateID   string     `json:"templateId"`
	Scope        string     `json:"scope"`
	Conditions   string     `json:"conditions"`
	HourlyRate   float64    `json:"hourlyRate"`
	CompanyName  string     `json:"companyName"`
	CompanyLogo  string     `gorm:"type:text" json:"companyLogo"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
