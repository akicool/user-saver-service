package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/akicool/user-saver-service/internal/proto/user"
	"github.com/akicool/user-saver-service/internal/user"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if (req.Name == "") || (req.Email == "") || (req.Password == "") {
		return &pb.CreateUserResponse{
			Message: "Все поля (name, email, password) обязательны",
			Status:  400,
		}, nil
	}

	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, req.Name, req.Email, req.Password)
	if err != nil {
		log.Println("Ошибка сохранения в БД:", err)
		return &pb.CreateUserResponse{
			Message: "Ошибка сохранения пользователя",
			Status:  500,
		}, err
	}

	return &pb.CreateUserResponse{
		Message: fmt.Sprintf("Пользователь %s успешно создан", req.Name),
		Status:  201,
	}, nil
}

func initDB() *sql.DB {
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
	return db
}

func main() {
	db := initDB()
	defer db.Close()

	internal.DB = db

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	fmt.Println("gRPC сервер запущен на порту 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
