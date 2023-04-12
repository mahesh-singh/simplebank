package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	//log.Printf(">> before: %v, %v", account1.Balance, account2.Balance)
	n := 5
	amount := int64(10)
	errs := make(chan error)
	txResults := make(chan TransferTxResult)
	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			txResults <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		res := <-txResults
		require.NotEmpty(t, res)

		//Check transfer
		require.NotEmpty(t, res.Transfer)
		require.Equal(t, account1.ID, res.Transfer.FromAccountID)
		require.Equal(t, account2.ID, res.Transfer.ToAccountID)
		require.Equal(t, amount, res.Transfer.Amount)

		//Entry
		require.Equal(t, res.FromEntry.AccountID, account1.ID)
		require.Equal(t, res.ToEntry.AccountID, account2.ID)
		require.Equal(t, res.Transfer.Amount, amount)

		//account
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance

		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)

		require.Equal(t, int64(0), diff2%amount)

	}

	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	//log.Printf(">> after: %v, %v", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedToAccount.Balance)

}
