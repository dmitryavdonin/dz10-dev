package repository

import (
	"order/internal/model"

	"gorm.io/gorm"
)

type Order interface {
	Create(order model.Order) (int, error)
	GetById(orderId int) (model.Order, error)
	GetAll(limit int, offset int) ([]model.Order, error)
	Delete(orderId int) error
}

type Repository struct {
	Order
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
