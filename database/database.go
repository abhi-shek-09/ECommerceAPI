package database

import (
	"database/sql"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB(){
	if err:= godotenv.Load(); err != nil {
		log.Fatal("Error loading environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	log.Println("Connected to Postgres successfully.")
}

func CloseDB(){
	if DB != nil {
		DB.Close()
		log.Println("Connection closed successfully.")
	}
}