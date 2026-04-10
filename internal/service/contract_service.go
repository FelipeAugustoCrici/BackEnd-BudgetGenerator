package service

import (
	"errors"
	"time"

	"budgetgen/internal/dto"
	"budgetgen/internal/model"
	"budgetgen/internal/repository"

	"github.com/google/uuid"
)

type ContractService struct {
	repo *repository.ContractRepository
}

func NewContractService(repo *repository.ContractRepository) *ContractService {
	return &ContractService{repo: repo}
}

func (s *ContractService) GetByBudgetID(budgetID, userID uuid.UUID) (*model.Contract, error) {
	return s.repo.GetByBudgetID(budgetID, userID)
}

func (s *ContractService) List(userID uuid.UUID, q dto.ContractFilterQuery) ([]model.Contract, int64, error) {
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	return s.repo.List(userID, q)
}

func (s *ContractService) GetByID(id, userID uuid.UUID) (*model.Contract, error) {
	return s.repo.GetByID(id, userID)
}

func (s *ContractService) Create(userID uuid.UUID, req dto.CreateContractRequest) (*model.Contract, error) {
	contract := &model.Contract{
		UserID:      userID,
		ClientID:    req.ClientID,
		BudgetID:    req.BudgetID,
		Value:       req.Value,
		Status:      "draft",
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		AutoRenew:   req.AutoRenew,
		Description: req.Description,
	}
	return contract, s.repo.Create(contract)
}

func (s *ContractService) Update(id, userID uuid.UUID, req dto.UpdateContractRequest) (*model.Contract, error) {
	contract, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	contract.BudgetID = req.BudgetID
	contract.Value = req.Value
	contract.StartDate = req.StartDate
	contract.EndDate = req.EndDate
	contract.AutoRenew = req.AutoRenew
	contract.Description = req.Description
	return contract, s.repo.Update(contract)
}

func (s *ContractService) Send(id, userID uuid.UUID) (*model.Contract, error) {
	contract, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if contract.Status != "draft" {
		return nil, errors.New("only draft contracts can be sent")
	}
	now := time.Now()
	contract.Status = "sent"
	contract.SentAt = &now
	if err := s.repo.Update(contract); err != nil {
		return nil, err
	}
	s.createEvent(contract.ID, "sent")
	return contract, nil
}

func (s *ContractService) View(id, userID uuid.UUID) (*model.Contract, error) {
	contract, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if contract.Status != "sent" {
		return nil, errors.New("contract must be in sent status to be viewed")
	}
	now := time.Now()
	contract.Status = "viewed"
	contract.ViewedAt = &now
	if err := s.repo.Update(contract); err != nil {
		return nil, err
	}
	s.createEvent(contract.ID, "viewed")
	return contract, nil
}

func (s *ContractService) Sign(id, userID uuid.UUID) (*model.Contract, error) {
	contract, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if contract.Status != "viewed" && contract.Status != "sent" {
		return nil, errors.New("contract must be sent or viewed to be signed")
	}
	now := time.Now()
	contract.Status = "signed"
	contract.SignedAt = &now
	if err := s.repo.Update(contract); err != nil {
		return nil, err
	}
	s.createEvent(contract.ID, "signed")
	return contract, nil
}

func (s *ContractService) Refuse(id, userID uuid.UUID) (*model.Contract, error) {
	contract, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if contract.Status == "signed" {
		return nil, errors.New("signed contracts cannot be refused")
	}
	contract.Status = "refused"
	if err := s.repo.Update(contract); err != nil {
		return nil, err
	}
	s.createEvent(contract.ID, "refused")
	return contract, nil
}

func (s *ContractService) ListEvents(contractID, userID uuid.UUID) ([]model.ContractEvent, error) {
	// verify ownership
	if _, err := s.repo.GetByID(contractID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListEvents(contractID)
}

func (s *ContractService) createEvent(contractID uuid.UUID, eventType string) {
	event := &model.ContractEvent{
		ContractID: contractID,
		Type:       eventType,
		Metadata:   "{}",
	}
	_ = s.repo.CreateEvent(event)
}
