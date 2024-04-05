package handler

import (
	"github.com/Just-Goo/Swift_Bank/service"
	"github.com/gin-gonic/gin"
)

type handlerImpl struct {
	service service.ServiceProvider
}

func NewHandlerImpl(s service.ServiceProvider) *handlerImpl  {
	return &handlerImpl{
		service: s,
	}
}

func (h *handlerImpl) SignUp(ctx *gin.Context)  {
	
}