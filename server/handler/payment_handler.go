package handler

import (
	"amartha-test/helper"
	"amartha-test/server/request"
	"amartha-test/usecase"
	"net/http"

	"github.com/go-chi/chi"
	validator "gopkg.in/go-playground/validator.v9"
)

// PaymentHandler ...
type PaymentHandler struct {
	Handler
}

// Execute ...
func (h *PaymentHandler) Execute(w http.ResponseWriter, r *http.Request) {
	req := request.MakePaymentRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	paymentUc := usecase.PaymentUC{ContractUC: h.ContractUC}
	err := paymentUc.MakePayment(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, nil)
}

// GetOutstanding ...
func (h *PaymentHandler) GetOutstanding(w http.ResponseWriter, r *http.Request) {
	loanId := chi.URLParam(r, "loanId")
	if loanId == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	paymentUc := usecase.PaymentUC{ContractUC: h.ContractUC}
	res, err := paymentUc.GetOutstanding(loanId)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			SendNotFound(w, err.Error())
			return
		}

		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}

// GetDelinquent ...
func (h *PaymentHandler) GetDelinquent(w http.ResponseWriter, r *http.Request) {
	loanId := chi.URLParam(r, "loanId")
	if loanId == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	paymentUc := usecase.PaymentUC{ContractUC: h.ContractUC}
	res, err := paymentUc.GetDelinquent(loanId)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			SendNotFound(w, err.Error())
			return
		}

		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}
