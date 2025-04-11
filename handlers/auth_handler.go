package handlers

import (
	"database/sql"
	"e-commerce/database"
	"e-commerce/models"
	"e-commerce/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func RegisterUser(w http.ResponseWriter, r *http.Request){
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
	}

	query := "INSERT INTO users (username, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id"
	err = database.DB.QueryRow(query, user.Username, user.Email, hashedPassword, user.IsAdmin).Scan(&user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){
			http.Error(w, "Email already exists in the database", http.StatusConflict)
		} else{
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{
		Token string `json:"token"`
	}{Token: token})
}

func LoginUser(w http.ResponseWriter, r *http.Request){
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var user models.User
	query := "SELECT id, password, is_admin FROM users WHERE email=$1"
	err := database.DB.QueryRow(query, creds.Email).Scan(&user.ID, &user.Password, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if !utils.VerifyPassword(user.Password, creds.Password){
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{
		Token string `json:"token"`
	}{Token: token})
}