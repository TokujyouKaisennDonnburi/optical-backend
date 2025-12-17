package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	calendarGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/gateway"
	calendarHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/handler"
	calendarCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	calendarQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	githubGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/github/gateway"
	githubHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/github/handler"
	githubCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	githubQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query"
	optionGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/option/gateway"
	optionHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/option/handler"
	optionQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/option/service/query"
	userGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/user/gateway"
	userHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	userCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	userQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	_ = godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://tokujyoukaisenndonnburi.github.io"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	jwtMiddleware := userHandler.NewUserAuthMiddleware()

	db := getPostgresDB()
	redisClient := GetRedisClient()
	minioClient := GetMinIOClient()

	// Migration
	MigrateMinio(minioClient)

	userRepository := userGateway.NewUserPsqlRepository(db)
	tokenRepository := userGateway.NewTokenRedisRepository(redisClient)
	stateRepository := githubGateway.NewStateRedisRepository(db, redisClient)
	optionRepository := optionGateway.NewOptionPsqlRepository(db)
	githubRepository := githubGateway.NewGithubApiRepository(db)
	githubQuery := githubQuery.NewGithubQuery(stateRepository, optionRepository, githubRepository)
	githubCommand := githubCommand.NewGithubCommand(tokenRepository, stateRepository, githubRepository)
	githubHandler := githubHandler.NewGithubHandler(githubQuery, githubCommand)
	userQuery := userQuery.NewUserQuery(userRepository)
	userCommand := userCommand.NewUserCommand(userRepository, tokenRepository)
	userHandler := userHandler.NewUserHttpHandler(userQuery, userCommand)
	optionQuery := optionQuery.NewOptionQuery(optionRepository)
	optionHandler := optionHandler.NewOptionHttpHandler(optionQuery)
	eventRepository := calendarGateway.NewEventPsqlRepository(db)
	calendarRepository := calendarGateway.NewCalendarPsqlRepository(db)
	imageRepository := calendarGateway.NewImagePsqlAndMinioRepository(db, minioClient, getBucketName())
	memberRepository := calendarGateway.NewMemberPsqlRepository(db)
	eventCommand := calendarCommand.NewEventCommand(eventRepository)
	eventQuery := calendarQuery.NewEventQuery(eventRepository)
	calendarCommand := calendarCommand.NewCalendarCommand(calendarRepository, optionRepository, imageRepository, memberRepository)
	calendarQuery := calendarQuery.NewCalendarQuery(calendarRepository)
	calendarHandler := calendarHandler.NewCalendarHttpHandler(eventCommand, eventQuery, calendarCommand, calendarQuery)

	// Unprotected Routes
	r.Group(func(r chi.Router) {
		// Users
		r.Post("/register", userHandler.Create)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)

		// Github
		r.Post("/github/apps/install", githubHandler.InstallToCalendar)
		r.Post("/github/oauth/link", githubHandler.LinkUser)
		r.Post("/github/oauth/create", githubHandler.CreateNewUserOauthState)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JWTAuthorization)

		// Users
		r.Get("/users/@me", userHandler.GetMe)
		r.Patch("/users/@me", userHandler.UpdateMe)

		// Github
		r.Post("/github/apps/state", githubHandler.CreateAppState)
		r.Post("/github/oauth/state", githubHandler.CreateOauthState)
		r.Post("/github/calendars/{calendarId}/review-requests", githubHandler.GetReviewRequests)
		r.Get("/github/calendars/{calendarId}/review-load-status", githubHandler.GetReviewLoadStatus)

		// Calendars
		r.Post("/calendars", calendarHandler.CreateCalendar)
		r.Post("/calendars/images", calendarHandler.UploadImage)
		r.Get("/calendars", calendarHandler.GetCalendars)

		// Members
		r.Post("/calendars/{calendarId}/members", calendarHandler.CreateMembers)
		r.Patch("/calendars/{calendarId}/members", calendarHandler.JoinMember)
		r.Delete("/calendars/{calendarId}/members", calendarHandler.RejectMember)

		// Events
		r.Post("/calendars/{calendarId}/events", calendarHandler.CreateEvent)
		r.Patch("/calendars/{calendarId}/events/{eventId}", calendarHandler.UpdateEvent)
		r.Get("/calendars/{calendarId}/events", calendarHandler.ListGetEvents)
		r.Get("/events/todays", calendarHandler.GetToday)
		r.Get("/events/months", calendarHandler.GetByMonth)

		// Options
		r.Get("/options", optionHandler.GetList)
	})

	// Start Serving
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}

func MigrateMinio(client *minio.Client) {
	exists, err := client.BucketExists(context.Background(), getBucketName())
	if err != nil {
		panic(err)
	}
	if exists {
		return
	}
	err = client.MakeBucket(context.Background(), getBucketName(), minio.MakeBucketOptions{
		Region: getMinioRegion(),
	})
	if err != nil {
		panic(err)
	}
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

func GetMinIOClient() *minio.Client {
	endpoint, ok := os.LookupEnv("MINIO_ENDPOINT")
	if !ok {
		panic("'MINIO_ENDPOINT' is not set'")
	}
	accessKeyId, ok := os.LookupEnv("MINIO_ACCESS_KEY_ID")
	if !ok {
		panic("'MINIO_ACCESS_KEY_ID' is not set'")
	}
	secretAccessKey, ok := os.LookupEnv("MINIO_SECRET_ACCESS_KEY")
	if !ok {
		panic("'MINIO_SECRET_ACCESS_KEY' is not set")
	}
	useSsl, ok := os.LookupEnv("MINIO_USE_SSL")
	if !ok {
		panic("'MINIO_USE_SSL' is not set")
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSsl == "1",
	})
	if err != nil {
		panic(err.Error())
	}
	return client
}

func getBucketName() string {
	bucketName, ok := os.LookupEnv("MINIO_IMAGE_BUCKET_NAME")
	if !ok {
		panic("'MINIO_IMAGE_BUCKET_NAME' is not set")
	}
	return bucketName
}

func getMinioRegion() string {
	bucketName, ok := os.LookupEnv("MINIO_REGION")
	if !ok {
		panic("'MINIO_REGION' is not set")
	}
	return bucketName
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
