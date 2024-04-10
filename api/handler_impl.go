package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/helpers"
	"github.com/Just-Goo/Swift_Bank/models"
	"github.com/Just-Goo/Swift_Bank/service"
	"github.com/Just-Goo/Swift_Bank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type handlerImpl struct {
	service    service.ServiceProvider
	tokenMaker token.Maker
	router     *gin.Engine
	config     config.Config
}

func newHandlerImpl(c config.Config, s service.ServiceProvider) (*handlerImpl, error) {
	tokenMaker, err := token.NewPastoMaker(c.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	handlerImpl := handlerImpl{
		service:    s,
		config:     c,
		tokenMaker: tokenMaker,
	}
	r := gin.New()

	// register currency validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", helpers.ValidCurrency)
	}

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

	return &handlerImpl, nil
}

func (h *handlerImpl) GetGin() *gin.Engine {
	return h.router
}

func (h *handlerImpl) GetTokenMaker() token.Maker {
	return h.tokenMaker
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

		v1.POST("/user", h.CreateUser)
		v1.POST("/users/login", h.LoginUser)

		authRoutes := v1.Group("/").Use(auth(h.tokenMaker))

		authRoutes.POST("/account", h.CreateAccount)
		authRoutes.GET("/account/:id", h.GetAccount)
		authRoutes.GET("/accounts", h.ListAccounts)
		authRoutes.PUT("/account/:id", h.UpdateAccount)
		authRoutes.DELETE("/account/:id", h.DeleteAccount)

		authRoutes.GET("/users/:id", h.GetUser)
		authRoutes.GET("/users", h.ListUsers)
		authRoutes.POST("/transfer", h.TransferMoney)
	}

}

func (h *handlerImpl) CreateAccount(ctx *gin.Context) {
	var req models.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	createdAccount, err := h.service.CreateAccount(ctx, req, authPayload.UserName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
			return
		}
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

	// authorization rule
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.UserName {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, h.errorResponse(err))
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	accounts, err := h.service.ListAccounts(ctx, authPayload.UserName, req.PageSize, ((req.PageID - 1) * req.PageSize))
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
		if err == pgx.ErrNoRows {
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
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.String(http.StatusOK, "account deleted")
}

func (h *handlerImpl) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}
	req.Password = hashedPassword

	createdUser, err := h.service.CreateUser(ctx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createdUser)
}

func (h *handlerImpl) LoginUser(ctx *gin.Context) {
	var req models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	user, err := h.service.LoginUser(ctx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusUnauthorized, h.errorResponse(err))
		return
	}

	accessToken, err := h.tokenMaker.CreateToken(req.UserName, h.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user":         user,
	})
}

func (h *handlerImpl) GetUser(ctx *gin.Context) {
	var req models.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	user, err := h.service.GetUser(ctx, req.UserName)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, h.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *handlerImpl) ListUsers(ctx *gin.Context) {
	var req models.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	users, err := h.service.ListUsers(ctx, req.PageSize, ((req.PageID - 1) * req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *handlerImpl) TransferMoney(ctx *gin.Context) {
	var req models.TransferMoneyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	fromAccount, valid := h.validAccount(ctx, req.FromAccountID, req.Currency)
	// check if the currencies match
	if !valid {
		return
	}

	// check if sender has enough money 
	if fromAccount.Balance < req.Amount {
		err := errors.New("sender does not have sufficient funds to transfer")
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return
	}

	// authorization rule
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.UserName {
		err := errors.New("sender account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, h.errorResponse(err))
		return
	}

	_, valid = h.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := models.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		Fee:           req.Fee,
	}

	createdAccount, err := h.service.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createdAccount)
}

func (h *handlerImpl) validAccount(ctx *gin.Context, accountID int64, currency string) (models.Account, bool) {
	account, err := h.service.GetAccount(ctx, accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, h.errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, h.errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, h.errorResponse(err))
		return account, false
	}

	return account, true
}

func (h *handlerImpl) errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
