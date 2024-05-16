package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zde37/Swift_Bank/helpers"
	mockedproviders "github.com/zde37/Swift_Bank/mock"
	"github.com/zde37/Swift_Bank/models"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {

	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockedproviders.NewMockRepositoryProvider(ctrl)
	// build stubs
	repo.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// create service
	service := NewService(repo)

	createdAccount, err := service.S.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)

	require.Equal(t, account, createdAccount)

}

func randomAccount() models.Account {
	return models.Account{
		ID:       helpers.RandomInt(1, 1000),
		Owner:    helpers.RandomOwner(),
		Balance:  float64(helpers.RandomMoney()),
		Currency: helpers.RandomCurrency(),
	}
}
