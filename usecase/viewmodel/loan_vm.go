package viewmodel

type LoanVM struct {
	ID           string  `json:"id"`
	Amount       float64 `json:"amount"`
	TotalAmount  float64 `json:"total_amount"`
	InterestRate float64 `json:"interest_rate"`
	Weeks        int     `json:"weeks"`
	CreatedAt    string  `json:"created_at"`
}
