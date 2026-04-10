package service

import (
	"budgetgen/internal/dto"
	"budgetgen/internal/model"
	"budgetgen/internal/repository"

	"github.com/google/uuid"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) List(userID uuid.UUID, q dto.ClientFilterQuery) ([]model.Client, int64, error) {
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	return s.repo.List(userID, q)
}

func (s *ClientService) GetByID(id, userID uuid.UUID) (*model.Client, error) {
	return s.repo.GetByID(id, userID)
}

func (s *ClientService) Create(userID uuid.UUID, req dto.CreateClientRequest) (*model.Client, error) {
	client := &model.Client{
		UserID:  userID,
		Name:    req.Name,
		Company: req.Company,
		Email:   req.Email,
		Phone:   req.Phone,
		Notes:   req.Notes,
	}
	return client, s.repo.Create(client)
}

func (s *ClientService) Update(id, userID uuid.UUID, req dto.UpdateClientRequest) (*model.Client, error) {
	client, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		client.Name = req.Name
	}
	client.Company = req.Company
	client.Email = req.Email
	client.Phone = req.Phone
	client.Notes = req.Notes
	return client, s.repo.Update(client)
}

func (s *ClientService) Delete(id, userID uuid.UUID) error {
	return s.repo.Delete(id, userID)
}
