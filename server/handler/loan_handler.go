package handler

import (
	"amartha-test/server/request"
	"amartha-test/usecase"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
)

// LoanHandler ...
type LoanHandler struct {
	Handler
}

// CreateLoan ...
func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	req := request.LoanRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	loanUc := usecase.LoanUC{ContractUC: h.ContractUC}
	res, err := loanUc.CreateLoan(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}
