package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iancardoso/pismo/internal/domain"
)

type TransactionsRepository struct {
	db *sql.DB
}

func NewTransactionsRepository(db *sql.DB) TransactionsRepository {
	return TransactionsRepository{db: db}
}

func (r TransactionsRepository) Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	result, err := r.db.ExecContext(
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

	transaction.ID = id
	return transaction, nil
}
