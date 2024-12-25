package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/akicool/user-saver-service/internal/proto/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	port := os.Getenv("GRPC_SERVER_PORT")
	if port == "" {
		log.Fatal("Порт не задан в .env файле")
	}

	serverAddress := fmt.Sprintf("localhost:%s", port)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите имя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		fmt.Println("Все поля обязательны для заполнения!")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}

	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка вызова метода CreateUser: %v", err)
	}

	fmt.Printf("Ответ от сервера: %s (статус: %d)\n", resp.Message, resp.Status)
}