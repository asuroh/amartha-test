package usecase

import (
	"amartha-test/model"
	"amartha-test/pkg/logruslogger"
	"amartha-test/server/request"
	"amartha-test/usecase/viewmodel"
	"database/sql"
	"fmt"
	"time"
)

const (
	uscaseLoanName = "LoanUC"
)

// LoanUC ...
type LoanUC struct {
	*ContractUC
	Tx *sql.Tx
}

func (uc LoanUC) CreateLoan(req *request.LoanRequest) (res viewmodel.LoanVM, err error) {
	var (
		funcName = "CreateLoan"
		ctx      = fmt.Sprintf("%s.%s", uscaseLoanName, funcName)
	)

	mLoan := model.NewLoanModel(uc.DB, uc.Tx)
	now := time.Now().Format(time.RFC3339)
	res = viewmodel.LoanVM{
		Amount:       req.Amount,
		TotalAmount:  calculateTotalAmount(req.Amount, req.InterestRate),
		InterestRate: req.InterestRate,
		Weeks:        req.Weeks,
		CreatedAt:    now,
	}
	res.ID, err = mLoan.CreateLoan(res)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	mPayment := model.NewPaymentModel(uc.DB, uc.Tx)
	err = mPayment.CreatePayment(genratePayment(res))
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	return res, nil
}

func genratePayment(loan viewmodel.LoanVM) []viewmodel.PaymentVM {
	var payments []viewmodel.PaymentVM
	for week := 1; week <= loan.Weeks; week++ {
		payments = append(payments, viewmodel.PaymentVM{
			LoanID: loan.ID,
			Week:   week,
			Amount: calculateAmountPerWeek(loan.TotalAmount, loan.Weeks),
		})
	}

	return payments
}

func calculateAmountPerWeek(totalAmount float64, weeks int) float64 {
	return totalAmount / float64(weeks)
}

func calculateTotalAmount(amount, interestRate float64) float64 {
	return amount * (1 + (interestRate / 100))
}
