package handler

import (
	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/service"
	"github.com/gin-gonic/gin"
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
