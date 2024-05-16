package service

import (
	"context"

	"github.com/zde37/Swift_Bank/models"
	"github.com/zde37/Swift_Bank/repository"
)

type ServiceProvider interface {
	CreateAccount(ctx context.Context, data models.CreateAccountRequest, username string) (models.Account, error)
	GetAccount(ctx context.Context, id int64) (models.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error)
	ListAccounts(ctx context.Context, naem string, limit, offset int32) ([]models.Account, error)
	UpdateAccount(ctx context.Context, id int64, balance float64) (models.Account, error)
	AddAccountBalance(ctx context.Context, id int64, balance float64) (models.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error)
	GetEntry(ctx context.Context, id int64) (models.Entry, error)
	ListEntries(ctx context.Context, accountID, limit, offset int64) ([]models.Entry, error)
	CreateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
	GetTransaction(ctx context.Context, id int64) (models.Transaction, error)
	ListTransactions(ctx context.Context, fromAccountID, toAccountID, limit, offset int64) ([]models.Transaction, error)
	TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error)
	CreateUser(ctx context.Context, data models.CreateUserRequest) (models.User, error)
	LoginUser(ctx context.Context, data models.LoginUserRequest) (models.User, error)
	GetUser(ctx context.Context, username string) (models.User, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]models.User, error)
}

type Service struct {
	S ServiceProvider
}

func NewService(repo repository.RepositoryProvider) *Service {
	return &Service{
		S: newServiceImpl(repo),
	}
}
