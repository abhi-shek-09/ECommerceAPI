package routes

import (
	"e-commerce/handlers"
	"e-commerce/middleware"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleWare)
	api.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	api.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	api.HandleFunc("/order", handlers.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", handlers.ViewOrders).Methods("GET")
	api.HandleFunc("/orders/{id:[0-9]+}", handlers.ViewOrderDetails).Methods("GET")
	api.HandleFunc("/orders/{id:[0-9]+}/cancel", handlers.CancelOrder).Methods("DELETE")
	api.HandleFunc("/cart", handlers.AddToCart).Methods("POST")
	api.HandleFunc("/cart", handlers.ViewCart).Methods("GET")
	api.HandleFunc("/cart/{product_id:[0-9]+}", handlers.RemoveFromCart).Methods("DELETE")

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminMiddleware)

	admin.HandleFunc("/products", handlers.AddProduct).Methods("POST")
	admin.HandleFunc("/products/{id:[0-9]+}", handlers.UpdateProduct).Methods("PUT")
	admin.HandleFunc("/products/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")
	admin.HandleFunc("/orders/{id:[0-9]+}/status", handlers.UpdateOrderStatus).Methods("PUT")

	// Payment routes
	api.HandleFunc("/create-payment-intent", handlers.CreatePaymentIntent).Methods("POST")
	api.HandleFunc("/webhook", handlers.HandleWebhook).Methods("POST")

	return router
}

// http://localhost:8080/api/admin/products  add a product
// http://localhost:8080/api/products   get all products
// http://localhost:8080/api/products/1   get all products
// http://localhost:8080/api/admin/products/1 update
// http://localhost:8080/api/admin/products/1 delete
