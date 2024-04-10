package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Just-Goo/Swift_Bank/helpers"
	mockedproviders "github.com/Just-Goo/Swift_Bank/mock"
	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserRequestMatcher struct {
	arg      models.CreateUserRequest
	password string
}

func (e eqCreateUserRequestMatcher) Matches(x interface{}) bool {
	arg, ok := x.(models.CreateUserRequest)
	if !ok {
		return false
	}

	err := helpers.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserRequestMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserRequest(a models.CreateUserRequest, p string) gomock.Matcher {
	return eqCreateUserRequestMatcher{arg: a, password: p}
}

func TestGetAccountAPI(t *testing.T) {

	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(service *mockedproviders.MockServiceProvider)
		checkResponse func(recorder *httptest.ResponseRecorder)
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server, err := NewHandler(testConfig, service)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/sb/api/v1/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.H.GetGin().ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateUserAPI(t *testing.T) {

	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(service *mockedproviders.MockServiceProvider)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.UserName,
				"fullname": user.FullName,
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				arg := models.CreateUserRequest{
					UserName: user.UserName,
					FullName: user.FullName,
					Email:    user.Email,
				}
				service.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserRequest(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWithUser(t, recorder.Body, user)
			},
		},
		{
			name: "INTERNAL ERROR",
			body: gin.H{
				"username": user.UserName,
				"fullname": user.FullName,
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, pgx.ErrTxClosed)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DUPLICATE USERNAME",
			body: gin.H{
				"username": user.UserName,
				"fullname": user.FullName,
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, &pgconn.PgError{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "INVALID USERNAME",
			body: gin.H{
				"username": "invalid-user#",
				"fullname": user.FullName,
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "INVALID EMAIL",
			body: gin.H{
				"username": user.UserName,
				"fullname": user.FullName,
				"password": password,
				"email":    "invalidemail.com",
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PASSWORD TOO SHORT",
			body: gin.H{
				"username": user.UserName,
				"fullname": user.FullName,
				"password": "123",
				"email":    user.Email,
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server, err := NewHandler(testConfig, service)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			// marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/sb/api/v1/user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.H.GetGin().ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
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

func randomUser(t *testing.T) (models.User, string) {
	return models.User{
		UserName: helpers.RandomOwner(),
		FullName: helpers.RandomOwner(),
		Email:    helpers.RandomEmail(),
	}, helpers.RandomString(6)
}

func requireBodyMatchWithUser(t *testing.T, body *bytes.Buffer, user models.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser models.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, gotUser, user)
}
