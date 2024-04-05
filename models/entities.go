package models

import "time"

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
