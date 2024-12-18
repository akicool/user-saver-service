package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Name, Email, Password string
}

var user User
var DB *sql.DB

func CreateUser(user User) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, user.Name, user.Email, user.Password)
	return err
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только метод POST поддерживается", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Все поля (name, email, password) обязательны", http.StatusBadRequest)
		return
	}
	
	if err := CreateUser(user); err != nil {
		http.Error(w, "Ошибка сохранения пользователя в базу данных", http.StatusInternalServerError)
		log.Println("Ошибка сохранения в БД:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Пользователь %s успешно создан", user.Name)
}