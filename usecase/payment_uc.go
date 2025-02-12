package usecase

import (
	"amartha-test/model"
	"amartha-test/pkg/logruslogger"
	"amartha-test/server/request"
	"amartha-test/usecase/viewmodel"
	"database/sql"
	"fmt"
)

const (
	uscasePaymentName = "PaymentUC"
)

// PaymentUC ...
type PaymentUC struct {
	*ContractUC
	Tx *sql.Tx
}

func (uc PaymentUC) MakePayment(req *request.MakePaymentRequest) (err error) {
	var (
		funcName = "MakePayment"
		ctx      = fmt.Sprintf("%s.%s", uscasePaymentName, funcName)
	)

	m := model.NewPaymentModel(uc.DB, uc.Tx)
	err = m.UpdatePayment(req.LoanId, req.Week)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return err
	}

	return nil
}

func (uc PaymentUC) GetOutstanding(loanId string) (res viewmodel.PaymentOutstandingVM, err error) {
	var (
		funcName = "GetOutstanding"
		ctx      = fmt.Sprintf("%s.%s", uscasePaymentName, funcName)
	)

	mPayment := model.NewPaymentModel(uc.DB, uc.Tx)
	totalPaid, err := mPayment.GetTotalPayment(loanId)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	payments, err := mPayment.GetPayments(loanId)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	for _, r := range payments {
		temp := viewmodel.PaymentOutstandingDetailVM{}
		ParsePayments(&r, &temp)
		res.Detail = append(res.Detail, temp)
	}

	mLoad := model.NewLoanModel(uc.DB, uc.Tx)
	totalAmount, err := mLoad.GetTotalAmount(loanId)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	res.Total = totalAmount - totalPaid

	return res, nil
}

func (uc PaymentUC) GetDelinquent(loanId string) (res viewmodel.DelinquentVM, err error) {
	var (
		funcName = "GetDelinquent"
		ctx      = fmt.Sprintf("%s.%s", uscasePaymentName, funcName)
	)

	mLoad := model.NewLoanModel(uc.DB, uc.Tx)
	weeksElapsed, err := mLoad.GetLoanWeeksElapsed(loanId)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	mPayment := model.NewPaymentModel(uc.DB, uc.Tx)
	res.IsDelinquent, err = mPayment.CheckDelinquent(loanId, weeksElapsed)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, funcName, uc.ReqID)
		return res, err
	}

	return res, nil
}

// ParsePayments ...
func ParsePayments(data *model.PaymentEntity, res *viewmodel.PaymentOutstandingDetailVM) {
	res.Amount = data.Amount
	res.Week = fmt.Sprintf("W%d", data.Week)
}
