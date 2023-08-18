package repository

import (
	"billing/internal/model"
	"time"

	"gorm.io/gorm"
)

type AccountPostgres struct {
	db *gorm.DB
}

func NewAccountPostgres(db *gorm.DB) *AccountPostgres {
	return &AccountPostgres{db: db}
}

func (r *AccountPostgres) Create(input model.BillingAccount) (int, error) {
	result := r.db.Create(&input)
	if result.Error != nil {
		return 0, result.Error
	}
	return input.ID, nil
}

func (r *AccountPostgres) GetById(id int) (model.BillingAccount, error) {
	var item model.BillingAccount
	result := r.db.First(&item, "user_id = ?", id)
	return item, result.Error
}

func (r *AccountPostgres) GetAll(limit int, offset int) ([]model.BillingAccount, error) {
	var items []model.BillingAccount
	result := r.db.Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}

func (r *AccountPostgres) Delete(id int) error {
	result := r.db.Delete(&model.BillingAccount{}, "user_id = ?", id)
	return result.Error
}

// update user balance or username
func (r *AccountPostgres) Update(id int, input model.BillingAccount) error {
	var updated model.BillingAccount
	result := r.db.First(&updated, "user_id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	now := time.Now()
	itemToUpdate := model.BillingAccount{
		UserId:     input.UserId,
		Balance:    input.Balance,
		CreatedAt:  updated.CreatedAt,
		ModifiedAt: now,
	}
	result = r.db.Model(&updated).Select("balance", "modified_at").Updates(itemToUpdate)
	return result.Error
}
