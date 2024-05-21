package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zde37/Swift_Bank/helpers"
	"github.com/zde37/Swift_Bank/models"
)

func createRandomAccount(t *testing.T) models.Account {
	user := createRandomUser(t)

	arg := models.Account{
		Owner:    user.UserName,
		Balance:  float64(helpers.RandomMoney()),
		Currency: helpers.RandomCurrency(),
	}

	account, err := testRepo.R.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createRandomEntry(t *testing.T, account *models.Account) models.Entry {
	arg := models.Entry{
		AccountID: account.ID,
		Amount:    float64(helpers.RandomMoney()),
	}

	entry, err := testRepo.R.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func createRandomTransaction(t *testing.T, account1, account2 *models.Account) models.Transaction {
	arg := models.Transaction{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Currency:      helpers.RandomCurrency(),
		Amount:        float64(helpers.RandomMoney()),
		Fee:           float64(helpers.RandomFee()),
		Description:   helpers.RandomString(20),
	}

	transaction, err := testRepo.R.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.FromAccountID, account1.ID)
	require.Equal(t, arg.ToAccountID, account2.ID)
	require.Equal(t, arg.Currency, transaction.Currency)
	require.Equal(t, arg.Amount, transaction.Amount)
	require.Equal(t, arg.Description, transaction.Description)
	require.Equal(t, arg.Fee, transaction.Fee)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)

	return transaction
}

func createRandomUser(t *testing.T) models.User {
	hashedPassword, err := helpers.HashPassword(helpers.RandomString(6))
	require.NoError(t, err)

	arg := models.User{
		UserName:       helpers.RandomOwner(),
		Email:          helpers.RandomEmail(),
		FullName:       helpers.RandomOwner(),
		HashedPassword: hashedPassword,
	}

	user, err := testRepo.R.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	// create account
	account1 := createRandomAccount(t)
	account2, err := testRepo.R.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// create user
	user1 := createRandomUser(t)
	user2, err := testRepo.R.GetUser(context.Background(), user1.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserFullName(t *testing.T) {
	user1 := createRandomUser(t)

	newFullName := helpers.RandomOwner()

	updatedUser, err := testRepo.R.UpdateUser(context.Background(), models.UpdateUserParams{
		UserName: user1.UserName,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, user1.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, user1.Email, updatedUser.Email)
	require.Equal(t, user1.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserEmail(t *testing.T) {
	user1 := createRandomUser(t)

	newEmail := helpers.RandomEmail()

	updatedUser, err := testRepo.R.UpdateUser(context.Background(), models.UpdateUserParams{
		UserName: user1.UserName,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, user1.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, user1.FullName, updatedUser.FullName)
	require.Equal(t, user1.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)

	newPassword, err := helpers.HashPassword(helpers.RandomString(10))
	require.NoError(t, err)

	updatedUser, err := testRepo.R.UpdateUser(context.Background(), models.UpdateUserParams{
		UserName: user1.UserName,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, user1.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newPassword, updatedUser.HashedPassword)
	require.Equal(t, user1.FullName, updatedUser.FullName)
	require.Equal(t, user1.Email, updatedUser.Email)
}

func TestUpdateAllUserDetails(t *testing.T) {
	user1 := createRandomUser(t)

	newFullName := helpers.RandomOwner()
	newEmail := helpers.RandomEmail()
	newPassword, err := helpers.HashPassword(helpers.RandomString(10))
	require.NoError(t, err)

	updatedUser, err := testRepo.R.UpdateUser(context.Background(), models.UpdateUserParams{
		UserName: user1.UserName,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, user1.FullName, updatedUser.FullName)
	require.NotEqual(t, user1.Email, updatedUser.Email)
	require.NotEqual(t, user1.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newPassword, updatedUser.HashedPassword)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	limit, offset := 5, 5

	users, err := testRepo.R.ListUsers(context.Background(), int32(limit), int32(offset))
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	newBalance := helpers.RandomMoney()

	account2, err := testRepo.R.UpdateAccount(context.Background(), account1.ID, float64(newBalance))
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, float64(newBalance), account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testRepo.R.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testRepo.R.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount models.Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	limit, offset := 5, 0

	accounts, err := testRepo.R.ListAccounts(context.Background(), lastAccount.Owner, int32(limit), int32(offset))
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestCreateEntry(t *testing.T) {
	account1 := createRandomAccount(t)
	createRandomEntry(t, &account1)
}

func TestGetEntry(t *testing.T) {
	account1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, &account1)

	entry2, err := testRepo.R.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1, entry2)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, &account)
	}

	limit, offset := 7, 3

	entries, err := testRepo.R.ListEntries(context.Background(), account.ID, int64(limit), int64(offset))
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 7)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}

func TestCreateTransaction(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransaction(t, &account1, &account2)
}

func TestGetTransaction(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transaction1 := createRandomTransaction(t, &account1, &account2)

	transaction2, err := testRepo.R.GetTransaction(context.Background(), transaction1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transaction2)

	require.Equal(t, transaction1, transaction2)
}

func TestListTransactions(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransaction(t, &account1, &account2)
		createRandomTransaction(t, &account1, &account2)
	}

	limit, offset := 10, 5
	transactions, err := testRepo.R.ListTransactions(context.Background(), account1.ID, account2.ID, int64(limit), int64(offset))
	require.NoError(t, err)
	require.NotEmpty(t, transactions)
	require.Len(t, transactions, 10)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
		require.Equal(t, transaction.FromAccountID, account1.ID)
		require.Equal(t, transaction.ToAccountID, account2.ID)
	}
}

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">== before: ", account1.Balance, account2.Balance)

	// run 'n' concurrent transfer transactions
	n := 5
	amount := float64(10)

	errs := make(chan error)
	results := make(chan models.TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testRepo.R.TransferTx(context.Background(), models.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      helpers.RandomCurrency(),
				Description:   helpers.RandomString(20),
				Fee:           float64(helpers.RandomFee()),
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transaction := result.Transaction
		require.NotEmpty(t, transaction)
		require.Equal(t, account1.ID, transaction.FromAccountID)
		require.Equal(t, account2.ID, transaction.ToAccountID)
		require.Equal(t, amount, transaction.Amount)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.CreatedAt)

		// check if the transaction was created in the database
		_, err = testRepo.R.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		// check the entries
		// sender entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// get sender entry from database
		_, err = testRepo.R.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// receiver entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// get receiver entry from database
		_, err = testRepo.R.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		// sending account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// receiving account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		fmt.Println(">== tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, int64(diff1)%int64(amount) == 0) // 1 * amount, 2 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testRepo.R.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testRepo.R.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">== after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, int64(account1.Balance)-int64(n)*int64(amount), int64(updatedAccount1.Balance))
	require.Equal(t, int64(account2.Balance)+int64(n)*int64(amount), int64(updatedAccount2.Balance))
}

func TestTransferTxDeadLock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">== before: ", account1.Balance, account2.Balance)

	// run 'n' concurrent transfer transactions
	n := 10
	amount := float64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testRepo.R.TransferTx(context.Background(), models.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
				Currency:      helpers.RandomCurrency(),
				Description:   helpers.RandomString(20),
				Fee:           float64(helpers.RandomFee()),
			})

			errs <- err
		}()
	}

	// check error
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testRepo.R.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testRepo.R.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">== after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
