package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/Swift_Bank/models"
)

type RepositoryProvider interface {
	CreateAccount(ctx context.Context, account models.Account) (models.Account, error)
	GetAccount(ctx context.Context, id int64) (models.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error)
	ListAccounts(ctx context.Context, name string, limit, offset int32) ([]models.Account, error)
	UpdateAccount(ctx context.Context, id int64, balance float64) (models.Account, error)
	AddAccountBalance(ctx context.Context, id int64, amount float64) (models.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error)
	GetEntry(ctx context.Context, id int64) (models.Entry, error)
	ListEntries(ctx context.Context, accountID, limit, offset int64) ([]models.Entry, error)
	CreateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
	GetTransaction(ctx context.Context, id int64) (models.Transaction, error)
	ListTransactions(ctx context.Context, fromAccountID, toAccountID, limit, offset int64) ([]models.Transaction, error)
	TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error)
	CreateUserTx(ctx context.Context, data models.CreateUserTxParams) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, userName string) (models.User, error)
	UpdateUser(ctx context.Context, user models.UpdateUserParams) (models.User, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]models.User, error)
	CreateSession(ctx context.Context, session models.Session) (models.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (models.Session, error)
}

type Repository struct {
	R RepositoryProvider
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		R: newRepositoryImpl(pool),
	}
}
