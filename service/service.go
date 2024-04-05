package service

import (
	"context"

	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/Just-Goo/Swift_Bank/repository"
)

type ServiceProvider interface {
	CreateAccount(ctx context.Context, data *models.Account) (*models.Account, error)
}

type Service struct {
	S ServiceProvider
}

func NewService(repo repository.RepositoryProvider) *Service {
	return &Service{
		S: NewServiceImpl(repo),
	}
}