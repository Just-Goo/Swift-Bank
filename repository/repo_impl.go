package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/Swift_Bank/models"
)

type repositoryImpl struct {
	pool *pgxpool.Pool
}

func newRepositoryImpl(p *pgxpool.Pool) *repositoryImpl {
	return &repositoryImpl{
		pool: p,
	}
}

func (r *repositoryImpl) CreateAccount(ctx context.Context, account models.Account) (models.Account, error) {
	var a models.Account
	query := `INSERT INTO accounts (owner, balance, currency) VALUES (@owner, @balance, @currency) RETURNING id, owner, balance, currency, created_at`
	args := pgx.NamedArgs{
		"owner":    account.Owner,
		"balance":  account.Balance,
		"currency": account.Currency,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency, &a.CreatedAt)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (r *repositoryImpl) GetAccount(ctx context.Context, id int64) (models.Account, error) {
	var account models.Account
	query := `SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repositoryImpl) GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error) {
	var account models.Account
	query := `SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = @id FOR NO KEY UPDATE`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repositoryImpl) ListAccounts(ctx context.Context, name string, limit, offset int32) ([]models.Account, error) {
	query := `SELECT id, owner, balance, currency, created_at FROM accounts WHERE owner = @name ORDER BY id LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"name":   name,
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []models.Account{}
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *repositoryImpl) UpdateAccount(ctx context.Context, id int64, balance float64) (models.Account, error) {
	query := "UPDATE accounts SET balance = @amount WHERE id = @id RETURNING id, owner, balance, currency, created_at"
	args := pgx.NamedArgs{
		"id":     id,
		"amount": balance,
	}

	var account models.Account

	err := r.pool.QueryRow(ctx, query, args).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repositoryImpl) AddAccountBalance(ctx context.Context, id int64, amount float64) (models.Account, error) {
	query := "UPDATE accounts SET balance = balance + @amount WHERE id = @id RETURNING id, owner, balance, currency, created_at"
	args := pgx.NamedArgs{
		"id":     id,
		"amount": amount,
	}

	var account models.Account

	err := r.pool.QueryRow(ctx, query, args).Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repositoryImpl) DeleteAccount(ctx context.Context, id int64) error {
	query := "DELETE FROM accounts WHERE id = @id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryImpl) CreateSession(ctx context.Context, session models.Session) (models.Session, error) {
	query := `INSERT INTO sessions (id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at) VALUES
			  (@id, @username, @token, @userAgent, @clientIP, @isBlocked, @expiresAt) RETURNING id, username, refresh_token, 
			  user_agent, client_ip, is_blocked, expires_at, created_at`
	args := pgx.NamedArgs{
		"id":        session.ID,
		"username":  session.UserName,
		"token":     session.RefreshToken,
		"userAgent": session.UserAgent,
		"clientIP":  session.ClientIp,
		"isBlocked": session.IsBlocked,
		"expiresAt": session.ExpiresAt,
	}

	var s models.Session
	err := r.pool.QueryRow(ctx, query, args).Scan(&s.ID, &s.UserName, &s.RefreshToken, &s.UserAgent, &s.ClientIp, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (r *repositoryImpl) GetSession(ctx context.Context, id uuid.UUID) (models.Session, error) {
	var s models.Session
	query := `SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&s.ID, &s.UserName, &s.RefreshToken, &s.UserAgent, &s.ClientIp, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (r *repositoryImpl) CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error) {
	query := `INSERT INTO entries (account_id, amount) VALUES (@accountID, @amount) RETURNING id, account_id, amount, created_at`
	args := pgx.NamedArgs{
		"accountID": entry.AccountID,
		"amount":    entry.Amount,
	}

	var e models.Entry
	err := r.pool.QueryRow(ctx, query, args).Scan(&e.ID, &e.AccountID, &e.Amount, &e.CreatedAt)
	if err != nil {
		return e, err
	}

	return e, nil
}

func (r *repositoryImpl) GetEntry(ctx context.Context, id int64) (models.Entry, error) {
	var entry models.Entry
	query := `SELECT id, account_id, amount, created_at FROM entries WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
	if err != nil {
		return entry, err
	}

	return entry, nil
}

func (r *repositoryImpl) ListEntries(ctx context.Context, accountID, limit, offset int64) ([]models.Entry, error) {
	query := `SELECT id, account_id, amount, created_at FROM entries WHERE account_id = @accountID ORDER BY id LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"accountID": accountID,
		"limit":     limit,
		"offset":    offset,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *repositoryImpl) CreateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	query := `INSERT INTO transactions (amount, fee, currency, description, to_account_id, from_account_id) VALUES 
	(@amount, @fee, @currency, @description, @toAccountID, @fromAccountID) RETURNING  id, from_account_id, to_account_id,
	 amount, fee, currency, description, created_at`
	args := pgx.NamedArgs{
		"amount":        transaction.Amount,
		"fee":           transaction.Fee,
		"currency":      transaction.Currency,
		"description":   transaction.Description,
		"fromAccountID": transaction.FromAccountID,
		"toAccountID":   transaction.ToAccountID,
	}

	var t models.Transaction
	err := r.pool.QueryRow(ctx, query, args).Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.Fee, &t.Currency, &t.Description, &t.CreatedAt)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (r *repositoryImpl) GetTransaction(ctx context.Context, id int64) (models.Transaction, error) {
	var transaction models.Transaction
	query := `SELECT id, from_account_id, to_account_id, amount, fee, currency, description, created_at FROM transactions WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&transaction.ID, &transaction.FromAccountID, &transaction.ToAccountID, &transaction.Amount, &transaction.Fee, &transaction.Currency, &transaction.Description, &transaction.CreatedAt)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repositoryImpl) ListTransactions(ctx context.Context, fromAccountID, toAccountID, limit, offset int64) ([]models.Transaction, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, fee, currency, description, created_at FROM transactions
				WHERE from_account_id = @fromAccountID OR to_account_id = @toAccountID ORDER BY id LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"fromAccountID": fromAccountID,
		"toAccountID":   toAccountID,
		"limit":         limit,
		"offset":        offset,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.FromAccountID, &transaction.ToAccountID, &transaction.Amount, &transaction.Fee, &transaction.Currency, &transaction.Description, &transaction.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *repositoryImpl) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var u models.User
	query := `INSERT INTO users (username, hashed_password, full_name, email) VALUES (@userName, @password, @fullName, @email) RETURNING username, hashed_password, full_name, email, password_changed_at, created_at`
	args := pgx.NamedArgs{
		"userName": user.UserName,
		"password": user.HashedPassword,
		"fullName": user.FullName,
		"email":    user.Email,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&u.UserName, &u.HashedPassword, &u.FullName, &u.Email, &u.PasswordChangedAt, &u.CreatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *repositoryImpl) GetUser(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := `SELECT username, hashed_password, full_name, email, password_changed_at, created_at FROM users WHERE username = @username`
	args := pgx.NamedArgs{
		"username": username,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&user.UserName, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repositoryImpl) ListUsers(ctx context.Context, limit, offset int32) ([]models.User, error) {
	query := `SELECT username, hashed_password, full_name, email, password_changed_at, created_at FROM users ORDER BY username LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserName, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repositoryImpl) UpdateUser(ctx context.Context, user models.UpdateUserParams) (models.User, error) {
	var u models.User
	query := `UPDATE users SET full_name = COALESCE(@newFullName, full_name),
			  hashed_password = COALESCE(@newPassword, hashed_password),
			  password_changed_at = COALESCE(@newPasswordChangedTime, password_changed_at),
			  email = COALESCE(@newEmail, email) WHERE username = @username  RETURNING 
			  username, hashed_password, full_name, email, password_changed_at, created_at`
	args := pgx.NamedArgs{
		"username":               user.UserName,
		"newFullName":            user.FullName,
		"newPassword":            user.HashedPassword,
		"newPasswordChangedTime": user.PasswordChangedAt,
		"newEmail":               user.Email,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&u.UserName, &u.HashedPassword, &u.FullName, &u.Email, &u.PasswordChangedAt, &u.CreatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *repositoryImpl) execTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (r *repositoryImpl) CreateUserTx(ctx context.Context, arg models.CreateUserTxParams) (models.User, error) {
	var result models.User

	err := r.execTx(ctx, func(tx pgx.Tx) error {
		var err error

		// create user
		query := `INSERT INTO users (username, hashed_password, full_name, email) VALUES (@userName, @password, @fullName, @email) RETURNING username, hashed_password, full_name, email, password_changed_at, created_at`
		args := pgx.NamedArgs{
			"userName": arg.User.UserName,
			"password": arg.User.HashedPassword,
			"fullName": arg.User.FullName,
			"email":    arg.User.Email,
		}

		err = tx.QueryRow(ctx, query, args).Scan(&result.UserName, &result.HashedPassword, &result.FullName, &result.Email, &result.PasswordChangedAt, &result.CreatedAt)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result) // callback function
	})

	return result, err
}

func (r *repositoryImpl) TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error) {
	var result models.TransferTxResult

	err := r.execTx(ctx, func(tx pgx.Tx) error {
		var err error
		// create transaction
		query := `INSERT INTO transactions (amount, fee, currency, description, to_account_id, from_account_id) VALUES 
					(@amount, @fee, @currency, @description, @toAccountID, @fromAccountID) RETURNING  id, from_account_id, to_account_id,
					 amount, fee, currency, description, created_at`
		args := pgx.NamedArgs{
			"amount":        arg.Amount,
			"fee":           arg.Fee,
			"currency":      arg.Currency,
			"description":   arg.Description,
			"fromAccountID": arg.FromAccountID,
			"toAccountID":   arg.ToAccountID,
		}

		err = tx.QueryRow(ctx, query, args).Scan(&result.Transaction.ID, &result.Transaction.FromAccountID, &result.Transaction.ToAccountID, &result.Transaction.Amount, &result.Transaction.Fee, &result.Transaction.Currency, &result.Transaction.Description, &result.Transaction.CreatedAt)
		if err != nil {
			return err
		}

		// create entry for sender account
		query2 := `INSERT INTO entries (account_id, amount) VALUES (@accountID, @amount) RETURNING id, account_id, amount, created_at`
		args2 := pgx.NamedArgs{
			"accountID": arg.FromAccountID,
			"amount":    -arg.Amount,
		}

		err = tx.QueryRow(ctx, query2, args2).Scan(&result.FromEntry.ID, &result.FromEntry.AccountID, &result.FromEntry.Amount, &result.FromEntry.CreatedAt)
		if err != nil {
			return err
		}

		// create entry for receiver account
		query3 := `INSERT INTO entries (account_id, amount) VALUES (@accountID, @amount) RETURNING id, account_id, amount, created_at`
		args3 := pgx.NamedArgs{
			"accountID": arg.ToAccountID,
			"amount":    arg.Amount,
		}

		err = tx.QueryRow(ctx, query3, args3).Scan(&result.ToEntry.ID, &result.ToEntry.AccountID, &result.ToEntry.Amount, &result.ToEntry.CreatedAt)
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			// update the sender's new balance first and then receiver's balance to avoid deadlock
			result.FromAccount, result.ToAccount, err = addMoney(ctx, tx, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			// update the receiver's new balance first and then sender's balance to avoid deadlock
			result.ToAccount, result.FromAccount, err = addMoney(ctx, tx, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	tx pgx.Tx,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amount2 float64,
) (account1 models.Account, account2 models.Account, err error) {
	// update account 1
	query := "UPDATE accounts SET balance = balance + @amount WHERE id = @id RETURNING id, owner, balance, currency, created_at"
	args := pgx.NamedArgs{
		"id":     accountID1,
		"amount": amount1,
	}

	err = tx.QueryRow(ctx, query, args).Scan(&account1.ID, &account1.Owner, &account1.Balance, &account1.Currency, &account1.CreatedAt)
	if err != nil {
		return account1, account2, err
	}

	// update account 2
	query2 := "UPDATE accounts SET balance = balance + @amount WHERE id = @id RETURNING id, owner, balance, currency, created_at"
	args2 := pgx.NamedArgs{
		"id":     accountID2,
		"amount": amount2,
	}

	err = tx.QueryRow(ctx, query2, args2).Scan(&account2.ID, &account2.Owner, &account2.Balance, &account2.Currency, &account2.CreatedAt)
	if err != nil {
		return account1, account2, err
	}

	return account1, account2, err
}
