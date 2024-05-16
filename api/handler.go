package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/service"
	"github.com/zde37/Swift_Bank/token"
)

type HandlerProvider interface {
	CreateAccount(ctx *gin.Context)
	GetAccount(ctx *gin.Context)
	ListAccounts(ctx *gin.Context)
	UpdateAccount(ctx *gin.Context)
	DeleteAccount(ctx *gin.Context)
	TransferMoney(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	ListUsers(ctx *gin.Context)
	GetGin() *gin.Engine
	StartServer(address string) error
	GetTokenMaker() token.Maker
}

type Handler struct {
	H HandlerProvider
}

func NewHandler(c config.Config, s service.ServiceProvider) (*Handler, error) {
	handler, err := newHandlerImpl(c, s)
	if err != nil {
		return nil, err
	}
	return &Handler{
		H: handler,
	}, nil
}
