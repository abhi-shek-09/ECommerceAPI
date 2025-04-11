package handlers

import (
	"e-commerce/database"
	"e-commerce/middleware"
	"e-commerce/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func AddProduct(w http.ResponseWriter, r *http.Request) { // ADMIN ONLY function
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	isAdmin, adminOk := r.Context().Value(middleware.IsAdminKey).(bool)

	if !ok || !adminOk || !isAdmin {
		http.Error(w, "Unauthorized: Only admins can create products", http.StatusForbidden)
		return
	}

	var product models.Products
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO products (name, description, price, stock) VALUES($1, $2, $3, $4) RETURNING id"
	err := database.DB.QueryRow(query, product.Name, product.Description, product.Price, product.Stock).Scan(&product.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := "SELECT id, name, description, price, stock FROM products"
	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Products
	for rows.Next() {
		var product models.Products
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			http.Error(w, "Error scanning products", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	productIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Products
	query := "SELECT id, name, description, price, stock FROM products WHERE id=$1"
	err = database.DB.QueryRow(query, productID).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	isAdmin, adminOk := r.Context().Value(middleware.IsAdminKey).(bool)

	if !ok || !adminOk || !isAdmin {
		http.Error(w, "Unauthorized: Only admins can create products", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	productIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Products
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5"
	_, err = database.DB.Exec(query, product.Name, product.Description, product.Price, product.Stock, productID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(int)              // User ID (not needed here)
	isAdmin, adminOk := r.Context().Value(middleware.IsAdminKey).(bool) // Extract admin flag

	if !ok || !adminOk || !isAdmin {
		http.Error(w, "Unauthorized: Only admins can create products", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	productIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM products WHERE id=$1"
	_, err = database.DB.Exec(query, productID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
