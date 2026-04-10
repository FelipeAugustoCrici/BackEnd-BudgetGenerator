package dto

import (
	"time"

	"github.com/google/uuid"
)

// Client DTOs
type CreateClientRequest struct {
	Name    string `json:"name" binding:"required"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Notes   string `json:"notes"`
}

type UpdateClientRequest struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Notes   string `json:"notes"`
}

// Contract DTOs
type CreateContractRequest struct {
	ClientID    uuid.UUID  `json:"client_id" binding:"required"`
	BudgetID    *uuid.UUID `json:"budget_id"`
	Value       float64    `json:"value"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	AutoRenew   bool       `json:"auto_renew"`
	Description string     `json:"description"`
}

type UpdateContractRequest struct {
	BudgetID    *uuid.UUID `json:"budget_id"`
	Value       float64    `json:"value"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	AutoRenew   bool       `json:"auto_renew"`
	Description string     `json:"description"`
}

// Pagination
type PaginationQuery struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=20"`
}

type ClientFilterQuery struct {
	PaginationQuery
	Search string `form:"search"`
}

type ContractFilterQuery struct {
	PaginationQuery
	ClientID string `form:"client_id"`
	Status   string `form:"status"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
