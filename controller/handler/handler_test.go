package handler

import (
	"bytes" 
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Just-Goo/Swift_Bank/helpers"
	mockedproviders "github.com/Just-Goo/Swift_Bank/mock"
	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {

	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(service *mockedproviders.MockServiceProvider)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWithAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NOT FOUND",
			accountID: account.ID,
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(models.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "INTERNAL ERROR",
			accountID: account.ID,
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(models.Account{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "INVALID ID",
			accountID: 0,
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mockedproviders.NewMockServiceProvider(ctrl)
			// build stubs
			tc.buildStubs(service)

			// create server
			server := NewHandler(service)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/sb/api/v1/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.H.GetGin().ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() models.Account {
	return models.Account{
		ID:       helpers.RandomInt(1, 1000),
		Owner:    helpers.RandomOwner(),
		Balance:  float64(helpers.RandomMoney()),
		Currency: helpers.RandomCurrency(),
	}
}

func requireBodyMatchWithAccount(t *testing.T, body *bytes.Buffer, account models.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount models.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)
}
