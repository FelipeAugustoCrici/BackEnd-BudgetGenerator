package repository

import (
	"budgetgen/internal/dto"
	"budgetgen/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) List(userID uuid.UUID, q dto.ClientFilterQuery) ([]model.Client, int64, error) {
	var clients []model.Client
	var total int64

	query := r.db.Model(&model.Client{}).Where("user_id = ?", userID)
	if q.Search != "" {
		like := "%" + q.Search + "%"
		query = query.Where("name ILIKE ? OR company ILIKE ? OR email ILIKE ?", like, like, like)
	}

	query.Count(&total)

	offset := (q.Page - 1) * q.Limit
	err := query.Offset(offset).Limit(q.Limit).Order("created_at DESC").Find(&clients).Error
	return clients, total, err
}

func (r *ClientRepository) GetByID(id, userID uuid.UUID) (*model.Client, error) {
	var client model.Client
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&client).Error
	return &client, err
}

func (r *ClientRepository) Create(client *model.Client) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) Update(client *model.Client) error {
	return r.db.Save(client).Error
}

func (r *ClientRepository) Delete(id, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Client{}).Error
}
