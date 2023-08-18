package service

import (
	"order/internal/model"
	"order/internal/repository"
)

type Order interface {
	Create(order model.Order) (int, error)
	GetById(orderId int) (model.Order, error)
	GetAll(limit int, offset int) ([]model.Order, error)
	Delete(orderId int) error
}

type Services struct {
	Order
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Order: NewOrderService(repos.Order),
	}
}
