package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/zde37/Swift_Bank/helpers"
	mockedproviders "github.com/zde37/Swift_Bank/mock"
	"github.com/zde37/Swift_Bank/models"
	"github.com/zde37/Swift_Bank/token"
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
	user, _ := randomUser()
	account := randomAccount(user.UserName)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWithAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UNAUTHORIZED USER",
			accountID: account.ID,
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NO AUTHORIZATION",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(service *mockedproviders.MockServiceProvider) {
				service.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
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

			tc.setupAuth(t, request, server.H.GetTokenMaker())
			server.H.GetGin().ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateUserAPI(t *testing.T) {

	user, password := randomUser()

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

func randomAccount(owner string) models.Account {
	return models.Account{
		ID:       helpers.RandomInt(1, 1000),
		Owner:    owner,
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

func randomUser() (models.User, string) {
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
