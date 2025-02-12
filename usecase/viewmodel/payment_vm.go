package viewmodel

type PaymentVM struct {
	ID     string  `json:"id"`
	LoanID string  `json:"loan_id"`
	Week   int     `json:"week"`
	Amount float64 `json:"amount"`
	PaidAt string  `json:"paid_at"`
	IsPaid bool    `json:"is_paid"`
}

type PaymentOutstandingVM struct {
	Total  float64                      `json:"total"`
	Detail []PaymentOutstandingDetailVM `json:"detail"`
}

type PaymentOutstandingDetailVM struct {
	Week   string  `json:"week"`
	Amount float64 `json:"amount"`
}
