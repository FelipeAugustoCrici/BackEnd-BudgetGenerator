package repository

import (
	"budgetgen/internal/dto"
	"budgetgen/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) List(userID uuid.UUID, q dto.ContractFilterQuery) ([]model.Contract, int64, error) {
	var contracts []model.Contract
	var total int64

	query := r.db.Model(&model.Contract{}).Preload("Client").Where("contracts.user_id = ?", userID)
	if q.ClientID != "" {
		query = query.Where("client_id = ?", q.ClientID)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}

	query.Count(&total)

	offset := (q.Page - 1) * q.Limit
	err := query.Offset(offset).Limit(q.Limit).Order("created_at DESC").Find(&contracts).Error
	return contracts, total, err
}

func (r *ContractRepository) GetByID(id, userID uuid.UUID) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.Preload("Client").Where("id = ? AND user_id = ?", id, userID).First(&contract).Error
	return &contract, err
}

func (r *ContractRepository) GetByBudgetID(budgetID, userID uuid.UUID) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.Preload("Client").Where("budget_id = ? AND user_id = ?", budgetID, userID).First(&contract).Error
	return &contract, err
}

func (r *ContractRepository) Create(contract *model.Contract) error {
	return r.db.Create(contract).Error
}

func (r *ContractRepository) Update(contract *model.Contract) error {
	return r.db.Save(contract).Error
}

func (r *ContractRepository) CreateEvent(event *model.ContractEvent) error {
	return r.db.Create(event).Error
}

func (r *ContractRepository) ListEvents(contractID uuid.UUID) ([]model.ContractEvent, error) {
	var events []model.ContractEvent
	err := r.db.Where("contract_id = ?", contractID).Order("created_at ASC").Find(&events).Error
	return events, err
}
