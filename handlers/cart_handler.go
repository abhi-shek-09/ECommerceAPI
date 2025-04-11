package handlers

import (
	"e-commerce/database"
	"e-commerce/middleware"
	"github.com/gorilla/mux"
	"e-commerce/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.(int)

	var cartItem models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id, created_at"
	err := database.DB.QueryRow(query, userID, cartItem.ProductID, cartItem.Quantity).Scan(&cartItem.ID, &cartItem.CreatedAt)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	cartItem.UserID = userID
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItem)
}

func ViewCart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.(int)

	query := "SELECT id, user_id, product_id, quantity, created_at FROM cart WHERE user_id=$1"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cartItems []models.Cart
	for rows.Next() {
		var cartItem models.Cart
		if err := rows.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt); err != nil {
			http.Error(w, "Error scanning cart items", http.StatusInternalServerError)
			return
		}
		cartItems = append(cartItems, cartItem)
	}

	if len(cartItems) == 0 {
		json.NewEncoder(w).Encode([]models.Cart{})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.(int)

	vars := mux.Vars(r)
	productIDStr, exists := vars["product_id"]
	if !exists {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM cart WHERE user_id=$1 AND product_id=$2"
	res, err := database.DB.Exec(query, userID, productID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Product not found in cart", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
