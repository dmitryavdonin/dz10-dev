package service

import (
	"user/internal/model"
	"user/internal/repository"
)

type User interface {
	Create(input model.User) (int, error)
	GetById(id int) (model.User, error)
	GetAll(limit int, offset int) ([]model.User, error)
	Delete(id int) error
	Update(id int, input model.User) error
}

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}
