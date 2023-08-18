package service

import (
	"billing/internal/model"
	"billing/internal/repository"
)

type Billing interface {
	Create(input model.Billing) (int, error)
	GetById(id int) (model.Billing, error)
	GetAll(limit int, offset int) ([]model.Billing, error)
	Delete(id int) error
	Update(id int, input model.Billing) error
}

type Service struct {
	Billing
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Billing: NewBillingService(repos.Billing),
	}
}
