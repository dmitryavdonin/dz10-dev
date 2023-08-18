package service

import (
	"billing/internal/model"
	"billing/internal/repository"
)

type BillingService struct {
	repo repository.Billing
}

func NewBillingService(repo repository.Billing) *BillingService {
	return &BillingService{repo: repo}
}

func (s *BillingService) Create(input model.Billing) (int, error) {
	return s.repo.Create(input)
}

func (s *BillingService) GetById(id int) (model.Billing, error) {
	return s.repo.GetById(id)
}

func (s *BillingService) GetAll(limit int, offset int) ([]model.Billing, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *BillingService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *BillingService) Update(id int, input model.Billing) error {
	return s.repo.Update(id, input)
}
