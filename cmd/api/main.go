package main

import (
	"fmt"
	"net/http"
	"os"

	userGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/user/gateway"
	userHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	userCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db := getPostgresDB()
	userRepository := userGateway.NewUserPsqlRepository(db)
	userCommand := userCommand.NewUserCommand(userRepository)
	userHandler := userHandler.NewUserHttpHandler(userCommand)

	// Users
	r.Post("/register", userHandler.Create)

	// Start Serving
	http.ListenAndServe(":8000", r)
}

// PostgresのDBに接続する
func getPostgresDB() *sqlx.DB {
	// #63 環境変数から読むように変える
	host := "127.0.0.1"
	port := "5433"
	sslMode := "disable"
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		panic("'POSTGRES_USER' is not set")
	}
	pass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		panic("'POSTGRES_PASSWORD' is not set")
	}
	dbName, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		panic("'POSTGRES_DB' is not set")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s database=%s user=%s password=%s sslmode=%s",
		host,
		port,
		dbName,
		user,
		pass,
		sslMode,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}
