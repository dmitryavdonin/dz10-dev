package service

import (
	"billing/internal/model"
	"billing/internal/repository"
)

type AccountService struct {
	repo repository.Account
}

func NewBillingService(repo repository.Account) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Create(input model.BillingAccount) (int, error) {
	return s.repo.Create(input)
}

func (s *AccountService) GetById(id int) (model.BillingAccount, error) {
	return s.repo.GetById(id)
}

func (s *AccountService) GetAll(limit int, offset int) ([]model.BillingAccount, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *AccountService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *AccountService) Update(id int, input model.BillingAccount) error {
	return s.repo.Update(id, input)
}
