package handlers

import (
	"database/sql"
	"e-commerce/database"
	"e-commerce/middleware"
	"e-commerce/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value(middleware.UserIDKey)
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    userID := user.(int)

    rows, err := database.DB.Query("SELECT product_id, quantity FROM cart WHERE user_id=$1", userID)
    if err != nil {
        http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var total float64
    var cartItems []models.Cart

    for rows.Next() {
        var cartItem models.Cart
        if err := rows.Scan(&cartItem.ProductID, &cartItem.Quantity); err != nil {
            http.Error(w, "Error scanning cart items", http.StatusInternalServerError)
            return
        }

        var price float64
        err := database.DB.QueryRow("SELECT price FROM products WHERE id=$1", cartItem.ProductID).Scan(&price)
        if err != nil {
            http.Error(w, "Failed to fetch product price", http.StatusInternalServerError)
            return
        }

        total += price * float64(cartItem.Quantity)
        cartItems = append(cartItems, cartItem)
    }

    if len(cartItems) == 0 {
        http.Error(w, "Cart is empty", http.StatusBadRequest)
        return
    }

    var order models.Orders
    query := "INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3) RETURNING id, created_at"
    err = database.DB.QueryRow(query, userID, total, "Not Paid").Scan(&order.ID, &order.CreatedAt)
    if err != nil {
        http.Error(w, "Failed to create order", http.StatusInternalServerError)
        return
    }

    // Clear the user's cart after the order is placed
    _, err = database.DB.Exec("DELETE FROM cart WHERE user_id=$1", userID)
    if err != nil {
        http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
        return
    }

    order.UserID = userID
    order.Total = total
    order.Status = "Not Paid"

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}


func ViewOrders(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.(int)

	query := "SELECT id, user_id, total, status, created_at FROM orders WHERE user_id=$1"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Orders
	for rows.Next() {
		var order models.Orders
		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		orders = []models.Orders{} // Ensure an empty array is returned instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func ViewOrderDetails(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := user.(int)
	vars := mux.Vars(r)
	orderIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.Orders
	query := "SELECT id, user_id, total, status, created_at FROM orders WHERE user_id=$1 AND id=$2"
	err = database.DB.QueryRow(query, userID, orderID).Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	isAdmin, adminOk := r.Context().Value(middleware.IsAdminKey).(bool)

	if !ok || !adminOk || !isAdmin {
		http.Error(w, "Unauthorized: Only admins can create products", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	orderIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// make a new struct to extract from request body
	var updateRequest struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	// map > arr, constant-time lookup
	validStatuses := map[string]bool{
		"Not Paid":  true,
		"Paid":      true,
		"Pending":   true,
		"Shipped":   true,
		"Delivered": true,
		"Cancelled": true,
	}

	if !validStatuses[updateRequest.Status] {
		http.Error(w, "Invalid order status", http.StatusBadRequest)
		return
	}

	query := "UPDATE orders SET status=$1 WHERE id=$2 RETURNING id, user_id, total, status, created_at"
	var updatedOrder models.Orders
	err = database.DB.QueryRow(query, updateRequest.Status, orderID).
		Scan(&updatedOrder.ID, &updatedOrder.UserID, &updatedOrder.Total, &updatedOrder.Status, &updatedOrder.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		}
		return
	}
		

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.(int)

	vars := mux.Vars(r)
	orderIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.Orders
	query := "SELECT id, user_id, total, status, created_at FROM orders WHERE user_id=$1 AND id=$2"
	err = database.DB.QueryRow(query, userID, orderID).Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		}
		return
	}
	

	if order.Status == "Shipped" || order.Status == "Delivered" {
		http.Error(w, fmt.Sprintf("Order cannot be cancelled as it is already %v", order.Status), http.StatusBadRequest)
		return
	}

	order.Status = "Cancelled"
	query = "UPDATE orders SET status=$1 WHERE id=$2 RETURNING id, user_id, total, status, created_at"
	err = database.DB.QueryRow(query, order.Status, orderID).
		Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		}
		return
	}
		

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
