package service

import (
	"context"
	"strings"

	"github.com/iancardoso/pismo/internal/domain"
	"github.com/iancardoso/pismo/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AccountsService struct {
	repo repository.AccountsRepository
}

func NewAccountsService(repo repository.AccountsRepository) AccountsService {
	return AccountsService{repo: repo}
}

func (s AccountsService) Create(ctx context.Context, documentNumber string) (domain.Account, error) {
	_, span := otel.Tracer("github.com/iancardoso/pismo/internal/service").Start(ctx, "accounts-service.create")
	defer span.End()

	documentNumber = strings.TrimSpace(documentNumber)
	span.SetAttributes(attribute.Bool("app.document_number.present", documentNumber != ""))
	if documentNumber == "" {
		span.SetStatus(codes.Error, ErrInvalidDocumentNumber.Error())
		return domain.Account{}, ErrInvalidDocumentNumber
	}

	exists, err := s.repo.ExistsByDocumentNumber(ctx, documentNumber)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return domain.Account{}, err
	}
	span.SetAttributes(attribute.Bool("app.account.document_exists", exists))
	if exists {
		span.SetStatus(codes.Error, ErrDuplicateDocument.Error())
		return domain.Account{}, ErrDuplicateDocument
	}

	account, err := s.repo.Create(ctx, documentNumber)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return domain.Account{}, err
	}
	span.SetAttributes(attribute.Int64("app.account.id", account.ID))
	return account, nil
}

func (s AccountsService) GetByID(ctx context.Context, id int64) (domain.Account, error) {
	_, span := otel.Tracer("github.com/iancardoso/pismo/internal/service").Start(ctx, "accounts-service.get-by-id")
	defer span.End()

	span.SetAttributes(attribute.Int64("app.account.id", id))
	if id <= 0 {
		span.SetStatus(codes.Error, ErrInvalidAccountID.Error())
		return domain.Account{}, ErrInvalidAccountID
	}

	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return domain.Account{}, err
	}
	return account, nil
}
