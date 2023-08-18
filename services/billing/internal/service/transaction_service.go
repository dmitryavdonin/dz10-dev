package service

import (
	"billing/internal/model"
	"billing/internal/repository"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Create(input model.BillingTransaction) (int, error) {
	return s.repo.Create(input)
}

func (s *TransactionService) GetById(id int) (model.BillingTransaction, error) {
	return s.repo.GetById(id)
}

func (s *TransactionService) GetAll(user_id int, limit int, offset int) ([]model.BillingTransaction, error) {
	return s.repo.GetAll(user_id, limit, offset)
}

func (s *TransactionService) Delete(id int) error {
	return s.repo.Delete(id)
}
