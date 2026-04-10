package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iancardoso/pismo/internal/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TransactionsHandler struct {
	service service.TransactionsService
}

type createTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func NewTransactionsHandler(service service.TransactionsService) TransactionsHandler {
	return TransactionsHandler{service: service}
}

// Create godoc
// @Summary Create a transaction
// @Accept json
// @Produce json
// @Param request body createTransactionRequest true "Create transaction request"
// @Success 201 {object} domain.Transaction
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /transactions [post]
func (h TransactionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("app.handler", "transactions.create"))

	var request createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.SetStatus(codes.Error, "invalid request body")
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	span.SetAttributes(
		attribute.Int64("app.account.id", request.AccountID),
		attribute.Int64("app.transaction.operation_type_id", request.OperationTypeID),
	)

	transaction, err := h.service.Create(r.Context(), request.AccountID, request.OperationTypeID, request.Amount)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		switch {
		case errors.Is(err, service.ErrInvalidAccountID), errors.Is(err, service.ErrInvalidOperationType), errors.Is(err, service.ErrInvalidAmount):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	span.SetAttributes(
		attribute.Int64("app.transaction.id", transaction.ID),
		attribute.Float64("app.transaction.amount", transaction.Amount),
	)
	writeJSON(w, http.StatusCreated, transaction)
}
