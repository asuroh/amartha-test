package request

// MakePaymentRequest...
type MakePaymentRequest struct {
	LoanId string `json:"loan_id"`
	Week   int    `json:"week"`
}
