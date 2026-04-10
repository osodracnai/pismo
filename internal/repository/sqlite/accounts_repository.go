package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iancardoso/pismo/internal/domain"
	"github.com/iancardoso/pismo/internal/service"
)

type AccountsRepository struct {
	db *sql.DB
}

func NewAccountsRepository(db *sql.DB) AccountsRepository {
	return AccountsRepository{db: db}
}

func (r AccountsRepository) Create(ctx context.Context, documentNumber string) (domain.Account, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO accounts (document_number) VALUES (?)`, documentNumber)
	if err != nil {
		return domain.Account{}, fmt.Errorf("create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Account{}, fmt.Errorf("read account id: %w", err)
	}

	return domain.Account{ID: id, DocumentNumber: documentNumber}, nil
}

func (r AccountsRepository) GetByID(ctx context.Context, id int64) (domain.Account, error) {
	var account domain.Account
	row := r.db.QueryRowContext(ctx, `SELECT id, document_number FROM accounts WHERE id = ?`, id)
	if err := row.Scan(&account.ID, &account.DocumentNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Account{}, service.ErrAccountNotFound
		}
		return domain.Account{}, fmt.Errorf("get account by id: %w", err)
	}
	return account, nil
}

func (r AccountsRepository) ExistsByDocumentNumber(ctx context.Context, documentNumber string) (bool, error) {
	var exists bool
	row := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM accounts WHERE document_number = ?)`, documentNumber)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("check document number: %w", err)
	}
	return exists, nil
}
