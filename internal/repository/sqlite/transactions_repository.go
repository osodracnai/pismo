package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iancardoso/pismo/internal/domain"
	"github.com/iancardoso/pismo/internal/service"
)

type TransactionsRepository struct {
	db *sql.DB
}

func NewTransactionsRepository(db *sql.DB) TransactionsRepository {
	return TransactionsRepository{db: db}
}

func (r TransactionsRepository) Apply(ctx context.Context, transaction domain.Transaction) (created domain.Transaction, err error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("open connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `BEGIN IMMEDIATE`); err != nil {
		return domain.Transaction{}, fmt.Errorf("begin immediate transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_, _ = conn.ExecContext(ctx, `ROLLBACK`)
		}
	}()

	var accountExists bool
	accountRow := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM accounts WHERE id = ?)`, transaction.AccountID)
	if err = accountRow.Scan(&accountExists); err != nil {
		return domain.Transaction{}, fmt.Errorf("check account exists: %w", err)
	}
	if !accountExists {
		return domain.Transaction{}, service.ErrAccountNotFound
	}

	var currentBalance float64
	balanceRow := conn.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE account_id = ?`, transaction.AccountID)
	if err = balanceRow.Scan(&currentBalance); err != nil {
		return domain.Transaction{}, fmt.Errorf("get account balance: %w", err)
	}
	if currentBalance+transaction.Amount < 0 {
		return domain.Transaction{}, service.ErrInsufficientFunds
	}

	result, err := conn.ExecContext(
		ctx,
		`INSERT INTO transactions (account_id, operation_type_id, amount, event_date) VALUES (?, ?, ?, ?)`,
		transaction.AccountID,
		transaction.OperationTypeID,
		transaction.Amount,
		transaction.EventDate.Format("2006-01-02T15:04:05.999999999Z07:00"),
	)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("read transaction id: %w", err)
	}

	if _, err = conn.ExecContext(ctx, `COMMIT`); err != nil {
		return domain.Transaction{}, fmt.Errorf("commit transaction: %w", err)
	}

	transaction.ID = id
	return transaction, nil
}

func (r TransactionsRepository) Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	return r.Apply(ctx, transaction)
}
