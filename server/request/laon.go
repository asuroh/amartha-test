package request

// LoanRequest ...
type LoanRequest struct {
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
	Weeks        int     `json:"weeks"`
}
