package repository

import (
	"billing/internal/model"

	"gorm.io/gorm"
)

type Billing interface {
	Create(item model.Billing) (int, error)
	GetById(id int) (model.Billing, error)
	GetAll(limit int, offset int) ([]model.Billing, error)
	Delete(id int) error
	Update(id int, item model.Billing) error
}

type Repository struct {
	Billing
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Billing: NewBillingPostgres(db),
	}
}
