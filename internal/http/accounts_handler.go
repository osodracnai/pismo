package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iancardoso/pismo/internal/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type AccountsHandler struct {
	service service.AccountsService
}

type createAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

func NewAccountsHandler(service service.AccountsService) AccountsHandler {
	return AccountsHandler{service: service}
}

// Create godoc
// @Summary Create an account
// @Accept json
// @Produce json
// @Param request body createAccountRequest true "Create account request"
// @Success 201 {object} domain.Account
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Router /accounts [post]
func (h AccountsHandler) Create(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("app.handler", "accounts.create"))

	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.SetStatus(codes.Error, "invalid request body")
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	span.SetAttributes(attribute.Bool("app.document_number.present", request.DocumentNumber != ""))

	account, err := h.service.Create(r.Context(), request.DocumentNumber)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		switch {
		case errors.Is(err, service.ErrInvalidDocumentNumber):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrDuplicateDocument):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	span.SetAttributes(attribute.Int64("app.account.id", account.ID))
	writeJSON(w, http.StatusCreated, account)
}

// GetByID godoc
// @Summary Get an account by ID
// @Produce json
// @Param accountId path int true "Account ID"
// @Success 200 {object} domain.Account
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /accounts/{accountId} [get]
func (h AccountsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("app.handler", "accounts.get_by_id"))

	accountID, err := strconv.ParseInt(chi.URLParam(r, "accountId"), 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, service.ErrInvalidAccountID.Error())
		writeError(w, http.StatusBadRequest, service.ErrInvalidAccountID.Error())
		return
	}
	span.SetAttributes(attribute.Int64("app.account.id", accountID))

	account, err := h.service.GetByID(r.Context(), accountID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		switch {
		case errors.Is(err, service.ErrInvalidAccountID):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, account)
}
