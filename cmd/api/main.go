package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	calendarGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/gateway"
	calendarHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/handler"
	calendarCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	optionGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/option/gateway"
	userGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/user/gateway"
	userHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	userCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	userQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	optionRepository := optionGateway.NewOptionPsqlRepository(db)
	eventRepository := calendarGateway.NewEventPsqlRepository(db)
	calendarRepository := calendarGateway.NewCalendarPsqlRepository(db)
	eventCommand := calendarCommand.NewEventCommand(eventRepository)
	calendarCommand := calendarCommand.NewCalendarCommand(calendarRepository, optionRepository)
	caledarHandler := calendarHandler.NewCalendarHttpHandler(eventCommand, calendarCommand)

	// Unprotected Routes
	r.Group(func(r chi.Router) {
		// Users
		r.Post("/register", userHandler.Create)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JWTAuthorization)

		// Users
		r.Get("/users/@me", userHandler.GetMe)

		// Calendars
		r.Post("/calendars", caledarHandler.CreateCalendar)

		// Events
		r.Post("/calendars/{calendarId}/events", caledarHandler.CreateEvent)
	})

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
