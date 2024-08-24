package domain

type Payment struct {
	PaymentID int `json:"payment_id"`
	Digital   bool
	COD       bool
	Credit    bool
}
