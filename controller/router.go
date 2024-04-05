package controller

import (
	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/controller/handler"
	"github.com/Just-Goo/Swift_Bank/controller/routes"
	"github.com/Just-Goo/Swift_Bank/repository"
	"github.com/Just-Goo/Swift_Bank/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(pool *pgxpool.Pool, config *config.Config)  {
	
	inProduction := false
	router := gin.New()

	router.Use(gin.Recovery())

	// set gin mode to release mode during production
	if inProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	if !inProduction {
		router.Use(gin.Logger())
	}

	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	handler := handler.NewHandler(service.S)
	routes.RegisterRoutes(router, handler)

	// default route
	router.GET("sb/api/v1", func(c *gin.Context) {
		c.Writer.Write([]byte("Swift Bank"))
	})

}