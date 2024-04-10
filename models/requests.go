package models

type CreateAccountRequest struct { 
	Currency string `json:"currency" binding:"required,currency"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateAccountRequest struct {
	Balance float64 `json:"balance" binding:"required"`
}

type CreateUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"` 
}

type GetUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
}

type TransferMoneyRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" binding:"required,min=1"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Fee           float64 `json:"fee" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required,currency"`
}
