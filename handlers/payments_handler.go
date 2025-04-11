package handlers

import (
	"e-commerce/database"
	"e-commerce/middleware"
	"e-commerce/models"
	"e-commerce/services"
	"encoding/json"
	"net/http"
)

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := user.(int)

	var req struct {
		OrderID int `json:"order_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON Request", http.StatusBadRequest)
		return
	}

	var order models.Orders
	err := database.DB.QueryRow("SELECT id, total FROM orders WHERE id=$1 AND user_id=$2", req.OrderID, userID).Scan(&order.ID, &order.Total)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	paymentIntent, err := services.CreatePaymentIntent(order.ID, int64(order.Total), "usd")
	if err != nil {
		http.Error(w, "Failed to create a payment intent", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO payments (user_id, order_id, amount, status) VALUES ($1, $2, $3, $4)"
	_, err = database.DB.Exec(query, userID, order.ID, order.Total, "pending")
	if err != nil {
		http.Error(w, "Failed to store payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"client_secret": paymentIntent.ClientSecret,
	})
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		OrderID int    `json:"order_id"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := services.UpdatePaymentStatus(payload.OrderID, payload.Status)
	if err != nil {
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
