package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var db *sql.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		dbUser, dbPassword, dbName, dbSSLMode, dbHost, dbPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	
	if err := db.Ping(); err != nil {
		log.Fatal("База данных недоступна:", err)
	}
	fmt.Println("Успешное подключение к базе данных")
}

func main() {
	initDB()

	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "user saver")
	})

	http.HandleFunc("/user-create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "user service")
	})

	fmt.Println("ListenAndServe localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}