package routes

import (
	"github.com/Just-Goo/Swift_Bank/controller/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handler.Handler) {
	r.POST("/signup", h.H.SignUp)
}
