package routes

import (
	"github.com/Just-Goo/Swift_Bank/controller/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handler.Handler) {
	v1 := r.Group("sb/api/v1")
	{
		// default route
		v1.GET("", func(c *gin.Context) {
			c.Writer.Write([]byte("Swift Bank"))
		})
	
		v1.POST("/account", h.H.SignUp)
		v1.GET("/account/:id", h.H.GetAccount)
		v1.GET("/accounts", h.H.ListAccounts)
		v1.PUT("/account/:id", h.H.UpdateAccount)
		v1.DELETE("/account/:id", h.H.DeleteAccount)
	}

}
