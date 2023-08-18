package repository

import (
	"billing/internal/model"

	"gorm.io/gorm"
)

type TransactionPostgres struct {
	db *gorm.DB
}

func NewTransactionPostgres(db *gorm.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) Create(input model.BillingTransaction) (int, error) {
	result := r.db.Create(&input)
	if result.Error != nil {
		return 0, result.Error
	}
	return input.ID, nil
}

func (r *TransactionPostgres) GetById(id int) (model.BillingTransaction, error) {
	var item model.BillingTransaction
	result := r.db.First(&item, "order_id = ?", id)
	return item, result.Error
}

func (r *TransactionPostgres) GetAll(user_id int, limit int, offset int) ([]model.BillingTransaction, error) {
	var items []model.BillingTransaction
	result := r.db.Limit(limit).Offset(offset).Find(&items, "user_id = ?", user_id)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}

func (r *TransactionPostgres) Delete(id int) error {
	result := r.db.Delete(&model.BillingTransaction{}, "order_id = ?", id)
	return result.Error
}
