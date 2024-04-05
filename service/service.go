package service

import (
	"context"

	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/Just-Goo/Swift_Bank/repository"
)

type ServiceProvider interface {
	CreateAccount(ctx context.Context, data *models.SignUpRequest) (models.Account, error)
	GetAccount(ctx context.Context, id int64) (*models.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (*models.Account, error)
	ListAccounts(ctx context.Context, limit, offset int32) ([]models.Account, error)
	UpdateAccount(ctx context.Context, id int64, balance float64) (*models.Account, error)
	AddAccountBalance(ctx context.Context, id int64, balance float64) (*models.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	CreateEntry(ctx context.Context, entry *models.Entry) (models.Entry, error)
	GetEntry(ctx context.Context, id int64) (*models.Entry, error)
	ListEntries(ctx context.Context, accountID, limit, offset int64) ([]models.Entry, error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (models.Transaction, error)
	GetTransaction(ctx context.Context, id int64) (*models.Transaction, error)
	ListTransactions(ctx context.Context, fromAccountID, toAccountID, limit, offset int64) ([]models.Transaction, error)
	TransferTx(ctx context.Context, arg *models.TransferTxParams) (models.TransferTxResult, error)
}

type Service struct {
	S ServiceProvider
}

func NewService(repo repository.RepositoryProvider) *Service {
	return &Service{
		S: NewServiceImpl(repo),
	}
}
