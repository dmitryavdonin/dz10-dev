package repository

import (
	"time"
	"user/internal/model"

	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(input model.User) (int, error) {
	result := r.db.Create(&input)
	if result.Error != nil {
		return 0, result.Error
	}
	return input.ID, nil
}

func (r *UserPostgres) GetById(id int) (model.User, error) {
	var item model.User
	result := r.db.First(&item, "id = ?", id)
	return item, result.Error
}

func (r *UserPostgres) GetAll(limit int, offset int) ([]model.User, error) {
	var items []model.User
	result := r.db.Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}

func (r *UserPostgres) Delete(id int) error {
	result := r.db.Delete(&model.User{}, "id = ?", id)
	return result.Error
}

// update user balance or username
func (r *UserPostgres) Update(id int, input model.User) error {
	var updated model.User
	result := r.db.First(&updated, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	now := time.Now()
	bookToUpdate := model.User{
		Username:   input.Username,
		Balance:    input.Balance,
		CreatedAt:  updated.CreatedAt,
		ModifiedAt: now,
	}
	result = r.db.Model(&updated).Select("username", "balance", "modified_at").Updates(bookToUpdate)
	return result.Error
}
