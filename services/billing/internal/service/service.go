package service

import (
	"billing/internal/model"
	"billing/internal/repository"
)

type Account interface {
	Create(input model.BillingAccount) (int, error)
	GetById(id int) (model.BillingAccount, error)
	GetAll(limit int, offset int) ([]model.BillingAccount, error)
	Delete(id int) error
	Update(id int, input model.BillingAccount) error
}

type Transaction interface {
	Create(input model.BillingTransaction) (int, error)
	GetById(id int) (model.BillingTransaction, error)
	GetAll(user_id int, limit int, offset int) ([]model.BillingTransaction, error)
	Delete(id int) error
}

type Services struct {
	Account
	Transaction
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Account:     NewBillingService(repos.Account),
		Transaction: NewTransactionService(repos.Transaction),
	}
}
