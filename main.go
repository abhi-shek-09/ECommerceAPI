package main

import (
	"e-commerce/database"
	"e-commerce/routes"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v78"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Fatal("STRIPE_SECRET_KEY not set in environment")
	}
	database.ConnectDB()
	defer database.CloseDB()

	router := routes.SetupRoutes()
	fmt.Println("Server started on :8080")
	log.Println(http.ListenAndServe(":8080", router))
}
