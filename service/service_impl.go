package service

import (
	"context"

	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/Just-Goo/Swift_Bank/repository"
)

type serviceImpl struct {
	repo repository.RepositoryProvider
}

func NewServiceImpl(r repository.RepositoryProvider) *serviceImpl {
	return &serviceImpl{
		repo: r,
	}
}

func (s *serviceImpl) CreateAccount(ctx context.Context, data *models.Account) (*models.Account, error)  {
	
	return nil, nil
}