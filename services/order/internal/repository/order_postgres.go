package repository

import (
	"order/internal/model"
	"time"

	"gorm.io/gorm"
)

type OrderPostgres struct {
	db *gorm.DB
}

func NewOrderPostgres(db *gorm.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) Create(order model.Order) (int, error) {
	result := r.db.Create(&order)
	if result.Error != nil {
		return 0, result.Error
	}

	return order.ID, nil
}

func (r *OrderPostgres) GetById(orderId int) (model.Order, error) {
	var order model.Order
	result := r.db.First(&order, "id = ?", orderId)
	return order, result.Error
}

func (r *OrderPostgres) GetAll(limit int, offset int) ([]model.Order, error) {
	var items []model.Order
	result := r.db.Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}

func (r *OrderPostgres) Delete(orderId int) error {
	result := r.db.Delete(&model.Order{}, "id = ?", orderId)
	return result.Error
}

func (r *OrderPostgres) Update(orderId int, input model.Order) error {
	var updatedItem model.Order
	result := r.db.First(&updatedItem, "id = ?", orderId)
	if result.Error != nil {
		return result.Error
	}
	now := time.Now()
	itemToUpdate := model.Order{
		Status:     input.Status,
		Reason:     input.Reason,
		ModifiedAt: now,
	}
	result = r.db.Model(&updatedItem).Select("status", "reason", "modified_at").Updates(itemToUpdate)
	return result.Error
}
