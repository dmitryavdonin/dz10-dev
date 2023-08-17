package repository

import (
	"user/internal/model"

	"gorm.io/gorm"
)

type User interface {
	Create(item model.User) (int, error)
	GetById(id int) (model.User, error)
	GetAll(limit int, offset int) ([]model.User, error)
	Delete(id int) error
	Update(id int, item model.User) error
}

type Repository struct {
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
