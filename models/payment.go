package models

import "time"

type Payments struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	OrderID       int       `json:"order_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}
