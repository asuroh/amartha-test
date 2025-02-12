package model

import (
	"amartha-test/usecase/viewmodel"
	"database/sql"
	"time"
)

// loanModel ...
type loanModel struct {
	DB *sql.DB
	Tx *sql.Tx
}

// ILoan ...
type ILoan interface {
	CreateLoan(body viewmodel.LoanVM) (string, error)
	GetTotalAmount(id string) (float64, error)
	GetLoanWeeksElapsed(id string) (weeksElapsed int, err error)
}

// LoanEntity ....
type LoanEntity struct {
	ID           string    `db:"id"`
	Amount       float64   `db:"amount"`
	TotalAmount  float64   `db:"total_amount"`
	InterestRate float64   `db:"interest_rate"`
	Weeks        int       `db:"weeks"`
	CreatedAt    time.Time `db:"created_at"`
}

// NewLoanModel ...
func NewLoanModel(db *sql.DB, tx *sql.Tx) ILoan {
	return &loanModel{DB: db, Tx: tx}
}

func (model loanModel) CreateLoan(body viewmodel.LoanVM) (string, error) {
	var id string
	sql := `INSERT INTO "loans" ("amount", "total_amount", "interest_rate", "weeks", "created_at") VALUES ($1, $2, $3, $4, $5) returning "id"`

	err := model.DB.QueryRow(sql, body.Amount, body.TotalAmount, body.InterestRate, body.Weeks, body.CreatedAt).Scan(&id)

	return id, err
}

func (model loanModel) GetTotalAmount(id string) (float64, error) {
	var totalAmount float64
	sql := `SELECT total_amount FROM loans WHERE id = $1`
	err := model.DB.QueryRow(sql, id).Scan(&totalAmount)
	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}

func (model loanModel) GetLoanWeeksElapsed(id string) (weeksElapsed int, err error) {
	var (
		createdAt time.Time
	)

	sql := `SELECT created_at FROM loans WHERE id = $1`
	err = model.DB.QueryRow(sql, id).Scan(&createdAt)
	if err != nil {
		return weeksElapsed, err
	}
	weeksElapsed = int(time.Since(createdAt).Hours() / (24 * 7))
	return weeksElapsed, nil
}
