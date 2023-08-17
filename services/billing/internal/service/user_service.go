package service

import (
	"user/internal/model"
	"user/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(input model.User) (int, error) {
	return s.repo.Create(input)
}

func (s *UserService) GetById(id int) (model.User, error) {
	return s.repo.GetById(id)
}

func (s *UserService) GetAll(limit int, offset int) ([]model.User, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) Update(id int, input model.User) error {
	return s.repo.Update(id, input)
}
