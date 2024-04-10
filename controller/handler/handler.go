package handler

import (
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

func NewHandler(s service.ServiceProvider) *Handler {
	return &Handler{
		H: newHandlerImpl(s),
	}
}
