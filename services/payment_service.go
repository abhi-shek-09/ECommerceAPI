package services

import (
	"e-commerce/database"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

func CreatePaymentIntent(userID int, amount int64, currency string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(currency),
	}
	
	return paymentintent.New(params)
}

func UpdatePaymentStatus(orderID int, status string) error {
	query := "UPDATE payments SET status=$1 WHERE order_id=$2"
	_, err := database.DB.Exec(query, status, orderID)
	return err
}
