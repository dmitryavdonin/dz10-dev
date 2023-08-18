package repository

import (
	"billing/internal/model"
	"time"

	"gorm.io/gorm"
)

type BillingPostgres struct {
	db *gorm.DB
}

func NewBillingPostgres(db *gorm.DB) *BillingPostgres {
	return &BillingPostgres{db: db}
}

func (r *BillingPostgres) Create(input model.Billing) (int, error) {
	result := r.db.Create(&input)
	if result.Error != nil {
		return 0, result.Error
	}
	return input.ID, nil
}

func (r *BillingPostgres) GetById(id int) (model.Billing, error) {
	var item model.Billing
	result := r.db.First(&item, "user_id = ?", id)
	return item, result.Error
}

func (r *BillingPostgres) GetAll(limit int, offset int) ([]model.Billing, error) {
	var items []model.Billing
	result := r.db.Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}

func (r *BillingPostgres) Delete(id int) error {
	result := r.db.Delete(&model.Billing{}, "user_id = ?", id)
	return result.Error
}

// update user balance or username
func (r *BillingPostgres) Update(id int, input model.Billing) error {
	var updated model.Billing
	result := r.db.First(&updated, "user_id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	now := time.Now()
	itemToUpdate := model.Billing{
		UserId:     input.UserId,
		Balance:    input.Balance,
		CreatedAt:  updated.CreatedAt,
		ModifiedAt: now,
	}
	result = r.db.Model(&updated).Select("balance", "modified_at").Updates(itemToUpdate)
	return result.Error
}
