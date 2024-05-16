package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    float64   `json:"amount"` // can be negative or positive
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Currency      string    `json:"currency"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"` // must be positive
	Fee           float64   `json:"fee"`    // must be positive
	CreatedAt     time.Time `json:"created_at"`
}

type User struct {
	UserName          string    `json:"username"`
	HashedPassword    string    `json:"-"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID    `json:"id"`
	UserName     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool    `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Description   string  `json:"description"`
	Fee           float64 `json:"fee"`
}

type TransferTxResult struct {
	Transaction Transaction `json:"transaction"`
	ToAccount   Account     `json:"to_account"`
	FromAccount Account     `json:"from_account"`
	ToEntry     Entry       `json:"to_entry"`
	FromEntry   Entry       `json:"from_entry"`
}
