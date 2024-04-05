package handler

import (
	"github.com/Just-Goo/Swift_Bank/service"
	"github.com/gin-gonic/gin"
)

type HandlerProvider interface {
	SignUp(ctx *gin.Context)
}

type Handler struct {
	H HandlerProvider
}

func NewHandler(s service.ServiceProvider) *Handler {
	return &Handler{
		H: NewHandlerImpl(s),
	}
}