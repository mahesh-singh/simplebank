package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	if store == nil {
		return errors.New("store is nil")
	}

	if store.db == nil {
		return errors.New("store.db is nil")
	}

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx.err: %v rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:to_account_id`
	Amount        int64 `json:amount`
}

type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:Entry`
	ToEntry     Entry    `json:Entry`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    arg.Amount})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount})
		if err != nil {
			return err
		}

		//This is to fix the deadlock. It will only update smaller account id first.
		//Ff two queries trying to update same account, it will prevent the deadlock
		//Otherwise, one query will update the account1 and waiting to get the lock on account 2,
		//where other query might update the account 2 and trying to acquire the lock on account 1.
		//This will create the deadlock.
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.FromAccountID, Amount: -arg.Amount})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.ToAccountID, Amount: arg.Amount})
			if err != nil {
				return err
			}
		} else {

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.ToAccountID, Amount: arg.Amount})
			if err != nil {
				return err
			}
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.FromAccountID, Amount: -arg.Amount})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
