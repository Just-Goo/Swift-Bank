package service

import (
	"context"

	"github.com/Just-Goo/Swift_Bank/helpers"
	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/Just-Goo/Swift_Bank/repository"
)

type serviceImpl struct {
	repo repository.RepositoryProvider
}

func newServiceImpl(r repository.RepositoryProvider) *serviceImpl {
	return &serviceImpl{
		repo: r,
	}
}

func (s *serviceImpl) CreateAccount(ctx context.Context, data models.CreateAccountRequest) (models.Account, error) {
	arg := models.Account{
		Owner:    data.Owner,
		Balance:  0,
		Currency: data.Currency,
	}

	return s.repo.CreateAccount(ctx, arg)
}

func (s *serviceImpl) GetAccount(ctx context.Context, id int64) (models.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *serviceImpl) GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error) {
	return s.repo.GetAccountForUpdate(ctx, id)
}

func (s *serviceImpl) ListAccounts(ctx context.Context, limit, offset int32) ([]models.Account, error) {
	return s.repo.ListAccounts(ctx, limit, offset)
}

func (s *serviceImpl) UpdateAccount(ctx context.Context, id int64, balance float64) (models.Account, error) {
	return s.repo.UpdateAccount(ctx, id, balance)
}

func (s *serviceImpl) AddAccountBalance(ctx context.Context, id int64, balance float64) (models.Account, error) {
	return s.repo.AddAccountBalance(ctx, id, balance)
}

func (s *serviceImpl) DeleteAccount(ctx context.Context, id int64) error {
	_, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteAccount(ctx, id)
}

func (s *serviceImpl) CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error) {

	return models.Entry{}, nil
}

func (s *serviceImpl) GetEntry(ctx context.Context, id int64) (models.Entry, error) {
	return s.repo.GetEntry(ctx, id)
}

func (s *serviceImpl) ListEntries(ctx context.Context, accountID, limit, offset int64) ([]models.Entry, error) {

	return nil, nil
}

func (s *serviceImpl) CreateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return models.Transaction{}, nil

}

func (s *serviceImpl) GetTransaction(ctx context.Context, id int64) (models.Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *serviceImpl) ListTransactions(ctx context.Context, fromAccountID, toAccountID, limit, offset int64) ([]models.Transaction, error) {

	return nil, nil
}

func (s *serviceImpl) CreateUser(ctx context.Context, data models.CreateUserRequest) (models.User, error) {
	hashedPassword, err := helpers.HashPassword(data.Password)
	if err != nil {
		return models.User{}, err
	}

	arg := models.User{
		UserName: data.UserName,
		FullName: data.FullName, 
		Email: data.Email,
		HashedPassword: hashedPassword,
	}

	return s.repo.CreateUser(ctx, arg)
}

func (s *serviceImpl) GetUser(ctx context.Context, username string) (models.User, error) {
	return s.repo.GetUser(ctx, username)
}

func (s *serviceImpl) ListUsers(ctx context.Context, limit, offset int32) ([]models.User, error) {
	return s.repo.ListUsers(ctx, limit, offset)
}

func (s *serviceImpl) TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error) {
	return s.repo.TransferTx(ctx, arg)
}
