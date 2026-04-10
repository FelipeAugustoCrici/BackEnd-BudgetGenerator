package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type QuoteSnapshot struct {
	ClientName   string     `json:"clientName"`
	Date         string     `json:"date"`
	Notes        string     `json:"notes"`
	Status       string     `json:"status"`
	Items        QuoteItems `json:"items"`
	Discount     float64    `json:"discount"`
	DiscountType string     `json:"discountType"`
	TemplateID   string     `json:"templateId"`
	Scope        string     `json:"scope"`
	Conditions   string     `json:"conditions"`
	HourlyRate   float64    `json:"hourlyRate"`
	CompanyName  string     `json:"companyName"`
}

func (s QuoteSnapshot) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *QuoteSnapshot) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, s)
}

type QuoteVersion struct {
	ID            string        `gorm:"primaryKey;type:uuid" json:"id"`
	QuoteID       string        `gorm:"type:uuid;not null;index" json:"quoteId"`
	UserID        string        `gorm:"type:uuid;not null;index" json:"userId"`
	VersionNumber int           `json:"versionNumber"`
	IsActive      bool          `gorm:"default:false" json:"isActive"`
	ChangeNote    string        `json:"changeNote"`
	Snapshot      QuoteSnapshot `gorm:"type:jsonb" json:"snapshot"`
	CreatedAt     time.Time     `json:"createdAt"`
}
