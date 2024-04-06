package handler

import (
	"net/http"

	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/Just-Goo/Swift_Bank/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type handlerImpl struct {
	service service.ServiceProvider
	router  *gin.Engine
}

func newHandlerImpl(s service.ServiceProvider) *handlerImpl {
	handlerImpl := handlerImpl{
		service: s,
	}
	r := gin.New()

	// inProduction := false
	// r.Use(gin.Recovery())

	// set gin mode to release mode during production
	// if inProduction {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// if !inProduction {
	r.Use(gin.Logger())
	// }
	
	handlerImpl.router = r
	handlerImpl.registerRoutes()

	return &handlerImpl
}

func (h *handlerImpl) GetGin() *gin.Engine {
	return h.router
}

func (h *handlerImpl) StartServer(address string) error {
	return h.router.Run(address)
}

func (h *handlerImpl) registerRoutes() {
	v1 := h.router.Group("sb/api/v1")
	{
		// default route
		v1.GET("", func(c *gin.Context) {
			c.Writer.Write([]byte("Swift Bank"))
		})

		v1.POST("/account", h.SignUp)
		v1.GET("/account/:id", h.GetAccount)
		v1.GET("/accounts", h.ListAccounts)
		v1.PUT("/account/:id", h.UpdateAccount)
		v1.DELETE("/account/:id", h.DeleteAccount)
	}

}

func (h *handlerImpl) SignUp(ctx *gin.Context) {
	var req models.SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	createdAccount, err := h.service.CreateAccount(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createdAccount)
}

func (h *handlerImpl) GetAccount(ctx *gin.Context) {
	var req models.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	account, err := h.service.GetAccount(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (h *handlerImpl) ListAccounts(ctx *gin.Context) {
	var req models.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	accounts, err := h.service.ListAccounts(ctx, req.PageSize, ((req.PageID - 1) * req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (h *handlerImpl) UpdateAccount(ctx *gin.Context) {
	var req models.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	var update models.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	account, err := h.service.AddAccountBalance(ctx, req.ID, update.Balance)
	if err != nil {
		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (h *handlerImpl) DeleteAccount(ctx *gin.Context) {
	var req models.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	err := h.service.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.String(http.StatusOK, "account deleted")
}

func (h *handlerImpl) errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
