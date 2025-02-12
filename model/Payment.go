package model

import (
	"amartha-test/usecase/viewmodel"
	"database/sql"
	"fmt"
	"time"
)

// PaymentModel ...
type PaymentModel struct {
	DB *sql.DB
	Tx *sql.Tx
}

func (model PaymentModel) scanRows(rows *sql.Rows) (d PaymentEntity, err error) {
	err = rows.Scan(
		&d.Week, &d.Amount,
	)

	return d, err
}

// IPayment ...
type IPayment interface {
	CreatePayment(payments []viewmodel.PaymentVM) error
	UpdatePayment(loanId string, week int) error
	GetTotalPayment(loanId string) (float64, error)
	GetPayments(loanId string) ([]PaymentEntity, error)
	CheckDelinquent(loanId string, weeksElapsed int) (bool, error)
}

// PaymentEntity ....
type PaymentEntity struct {
	ID     string    `db:"id"`
	LoanID string    `db:"loan_id"`
	Week   int       `db:"week"`
	Amount float64   `db:"amount"`
	PaidAt time.Time `db:"paid_at"`
	IsPaid bool      `db:"is_paid"`
}

// NewPaymentModel ...
func NewPaymentModel(db *sql.DB, tx *sql.Tx) IPayment {
	return &PaymentModel{DB: db, Tx: tx}
}

func (model PaymentModel) CreatePayment(payments []viewmodel.PaymentVM) error {
	sql := `INSERT INTO "payments" ("loan_id", "amount", "week") VALUES `
	values := []interface{}{}
	for i, payment := range payments {
		if i > 0 {
			sql += ", "
		}
		sql += fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
		values = append(values, payment.LoanID, payment.Amount, payment.Week)
	}
	_, err := model.DB.Exec(sql, values...)

	return err
}

func (model PaymentModel) UpdatePayment(loanId string, week int) (err error) {
	sql := `UPDATE "payments" SET "is_paid" = TRUE, "paid_at" = NOW() WHERE "loan_id" = $1 AND week = $2`
	if model.Tx != nil {
		_, err = model.Tx.Exec(sql, loanId, week)
		return err
	} else {
		_, err = model.DB.Exec(sql, loanId, week)
	}

	return err
}

func (model PaymentModel) GetTotalPayment(loanId string) (float64, error) {
	var totalPayment float64
	sql := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE loan_id = $1 AND is_paid = TRUE`
	err := model.DB.QueryRow(sql, loanId).Scan(&totalPayment)
	if err != nil {
		return 0, err
	}

	return totalPayment, err
}

func (model PaymentModel) GetPayments(loanId string) (res []PaymentEntity, err error) {
	sql := `SELECT week, amount FROM payments WHERE loan_id=$1 AND is_paid=FALSE ORDER BY week ASC`
	rows, err := model.DB.Query(sql, loanId)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		d, err := model.scanRows(rows)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}
	err = rows.Err()
	if err != nil {
		return res, err
	}

	return res, err
}

func (model PaymentModel) CheckDelinquent(loanId string, weeksElapsed int) (bool, error) {
	var missedPayments int
	sql := `SELECT COUNT(*) FROM payments WHERE loan_id = $1 AND is_paid = FALSE AND week <= $2`
	err := model.DB.QueryRow(sql, loanId, weeksElapsed).Scan(&missedPayments)
	if err != nil {
		return false, err
	}

	return missedPayments > 2, nil
}
