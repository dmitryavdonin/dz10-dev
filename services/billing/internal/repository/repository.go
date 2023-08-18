package repository

import (
	"billing/internal/model"

	"gorm.io/gorm"
)

type Account interface {
	Create(item model.BillingAccount) (int, error)
	GetById(id int) (model.BillingAccount, error)
	GetAll(limit int, offset int) ([]model.BillingAccount, error)
	Delete(id int) error
	Update(id int, item model.BillingAccount) error
}

type Transaction interface {
	Create(item model.BillingTransaction) (int, error)
	GetById(id int) (model.BillingTransaction, error)
	GetAll(user_id int, limit int, offset int) ([]model.BillingTransaction, error)
	Delete(id int) error
}

type Repository struct {
	Account
	Transaction
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Account:     NewAccountPostgres(db),
		Transaction: NewTransactionPostgres(db),
	}
}
