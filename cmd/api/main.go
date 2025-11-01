package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	userGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/user/gateway"
	userHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	userCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	userQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	jwtMiddleware := userHandler.NewUserAuthMiddleware()

	db := getPostgresDB()
	redisClient := GetRedisClient()
	userRepository := userGateway.NewUserPsqlRepository(db)
	tokenRepository := userGateway.NewTokenRedisRepository(redisClient)
	userQuery := userQuery.NewUserQuery(userRepository)
	userCommand := userCommand.NewUserCommand(userRepository, tokenRepository)
	userHandler := userHandler.NewUserHttpHandler(userQuery, userCommand)

	// Unprotected Routes
	r.Group(func(r chi.Router) {
		// Users
		r.Post("/register", userHandler.Create)
		r.Post("/login", userHandler.Login)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JWTAuthorization)

		r.Get("/users/@me", userHandler.GetMe)
	})

	// Start Serving
	http.ListenAndServe(":8000", r)
}

// PostgresのDBに接続する
func getPostgresDB() *sqlx.DB {
	// #63 環境変数から読むように変える
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		panic("'POSTGRES_HOST' is not set")
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		panic("'POSTGRES_PORT' is not set")
	}
	sslMode, ok := os.LookupEnv("POSTGRES_SSLMODE")
	if !ok {
		panic("'POSTGRES_SSLMODE' is not set")
	}
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

func GetRedisClient() *redis.Client {
	endpoint, ok := os.LookupEnv("REDIS_ADDRESS")
	if !ok {
		panic("'REDIS_ADDRESS' is not set")
	}
	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		panic("'REDIS_PASSWORD' is not set")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     endpoint,
		Password: password,
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("redis connection failed")
	}
	return client
}
