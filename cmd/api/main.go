package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	migrate "github.com/rubenv/sql-migrate"

	agentGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/agent/gateway"
	agentTransactor "github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	agentHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/agent/handler"
	agentCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/command"
	agentQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/query"
	calendarGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/gateway"
	calendarHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/handler"
	calendarCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	calendarQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	githubGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/github/gateway"
	githubHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/github/handler"
	githubCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	githubQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query"
	noticeRepository "github.com/TokujouKaisenDonburi/optical-backend/internal/notice/gateway"
	noticeHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/notice/handler"
	noticeQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/notice/service/query"
	optionGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/option/gateway"
	optionHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/option/handler"
	optionQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/option/service/query"
	schedulerGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/gateway"
	schedulerHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/handler"
	schedulerCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/command"
	schedulerQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query"
	todoGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/gateway"
	todoHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/handler"
	todoCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/command"
	todoQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/query"
	userGateway "github.com/TokujouKaisenDonburi/optical-backend/internal/user/gateway"
	userHandler "github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	userCommand "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	userQuery "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/logs"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/transact"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"google.golang.org/genai"
	"gopkg.in/mail.v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func main() {
	_ = godotenv.Load()

	// Logger configurations
	reportCaller := os.Getenv("LOGGER_REPORT_CALLER") == "1"
	logrus.SetReportCaller(reportCaller)
	if os.Getenv("LOGGER_JSON_FORMAT") == "1" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: false,
		})
	} else {
		formatter := &prefixed.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
		formatter.SetColorScheme(&prefixed.ColorScheme{
			TimestampStyle:  "white",
			ErrorLevelStyle: "red+b",
			FatalLevelStyle: "red+bu",
		})
		logrus.SetFormatter(formatter)
	}
	if level, err := logrus.ParseLevel(os.Getenv("LOGGER_LEVEL")); err == nil {
		logrus.SetLevel(level)
		logrus.WithField("level", level.String()).Info("loglevel updated")
	}

	r := chi.NewRouter()
	jwtMiddleware := userHandler.NewUserAuthMiddleware()

	db := getPostgresDB()
	dialer := GetDialer()
	redisClient := GetRedisClient()
	minioClient := GetMinIOClient()
	// genaiClient := GetGenAIClient()
	openRouter := GetOpenRouter()

	// Migration
	if os.Getenv("RUNTIME_MIGRATION") == "1" {
		migrations := &migrate.FileMigrationSource{
			Dir: "db/migrate",
		}
		n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			panic(err.Error())
		} else {
			logrus.WithField("applied", n).Info("migration executed")
		}
		MigrateMinio(minioClient)
	}
	bucketName := getBucketName()
	userRepository := userGateway.NewUserPsqlRepository(db)
	avatarRepository := userGateway.NewAvatarPsqlAndMinioRepository(db, minioClient, bucketName)
	tokenRepository := userGateway.NewTokenRedisRepository(redisClient)
	googleRepository := userGateway.NewGooglePsqlAndApiRepository(
		db,
		getGoogleOauthClientId(),
		getGoogleOauthClientSecret(),
		getGoogleOauthRedirectUri(),
	)
	googleOauthStateRepository := userGateway.NewGoogleOauthStateRedisRepository(
		getGoogleOauthClientId(),
		getGoogleOauthClientSecret(),
		getGoogleOauthRedirectUri(),
		redisClient,
	)
	redisEncryptionKey := getRedisEncryptionKey()
	installationIdEncryptionKey := getInstallationIdEncryptionKey()
	transactionProvider := transact.NewPsqlTransactionProvider(db)
	agentTransactor := agentTransactor.NewTransactionProvider(db)
	stateRepository := githubGateway.NewStateRedisRepository(db, redisClient, redisEncryptionKey)
	optionRepository := optionGateway.NewOptionPsqlRepository(db)
	githubRepository := githubGateway.NewGithubApiRepository(db, installationIdEncryptionKey)
	gmailRepository := calendarGateway.NewGmailRepository(dialer, getEmailAddress())
	eventRepository := calendarGateway.NewEventPsqlRepository(db)
	optionAgentRepository := agentGateway.NewOptionAgentOpenRouterRepository(openRouter)
	agentQueryRepository := agentGateway.NewAgentQueryPsqlRepository(db)
	calendarRepository := calendarGateway.NewCalendarPsqlRepository(db)
	imageRepository := calendarGateway.NewImagePsqlAndMinioRepository(db, minioClient, bucketName)
	memberRepository := calendarGateway.NewMemberPsqlRepository(db)
	agentCommandRepository := agentGateway.NewAgentCommandPsqlRepository(db)
	agentQuery := agentQuery.NewAgentQuery(optionRepository, eventRepository, optionAgentRepository)
	agentCommand := agentCommand.NewAgentCommand(openRouter, agentTransactor, agentQueryRepository, agentCommandRepository)
	agentHandler := agentHandler.NewAgentHandler(agentQuery, agentCommand)
	githubQuery := githubQuery.NewGithubQuery(stateRepository, optionRepository, githubRepository)
	githubCommand := githubCommand.NewGithubCommand(tokenRepository, stateRepository, githubRepository)
	githubHandler := githubHandler.NewGithubHandler(githubQuery, githubCommand)
	userQuery := userQuery.NewUserQuery(userRepository)
	userCommand := userCommand.NewUserCommand(
		transactionProvider,
		userRepository,
		tokenRepository,
		avatarRepository,
		googleRepository,
		calendarRepository,
		googleOauthStateRepository,
	)
	userHandler := userHandler.NewUserHttpHandler(userQuery, userCommand)
	optionQuery := optionQuery.NewOptionQuery(optionRepository)
	optionHandler := optionHandler.NewOptionHttpHandler(optionQuery)
	eventCommand := calendarCommand.NewEventCommand(transactionProvider, eventRepository, calendarRepository)
	eventQuery := calendarQuery.NewEventQuery(eventRepository)
	calendarCommand := calendarCommand.NewCalendarCommand(transactionProvider, calendarRepository, optionRepository, imageRepository, memberRepository, gmailRepository)
	memberQuery := calendarQuery.NewMemberQuery(memberRepository)
	calendarQuery := calendarQuery.NewCalendarQuery(calendarRepository)
	calendarHandler := calendarHandler.NewCalendarHttpHandler(eventCommand, eventQuery, calendarCommand, calendarQuery, memberQuery)
	noticeRepository := noticeRepository.NewNoticePsqlRepository(db)
	noticeQueryService := noticeQuery.NewNoticeQuery(noticeRepository)
	noticeHttpHandler := noticeHandler.NewNoticeHttpHandler(noticeQueryService)
	schedulerRepository := schedulerGateway.NewSchedulerPsqlRepository(db)
	schedulerCommand := schedulerCommand.NewSchedulerCommand(schedulerRepository, optionRepository)
	schedulerQuery := schedulerQuery.NewSchedulerQuery(schedulerRepository, optionRepository)
	schedulerHandler := schedulerHandler.NewSchedulerHttpHandler(schedulerCommand, schedulerQuery)

	r.Use(logs.HttpLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://tokujyoukaisenndonnburi.github.io", "https://opti-cal.org"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	todoRepository := todoGateway.NewTodoPsqlRepository(db)
	todoQuery := todoQuery.NewTodoQuery(todoRepository)
	todoCommand := todoCommand.NewTodoCommand(todoRepository)
	todoHandler := todoHandler.NewTodoHttpHandler(todoQuery, todoCommand)

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

		// Google
		r.Post("/google/oauth/user", userHandler.CreateGoogleUser)
		r.Post("/google/oauth/state", userHandler.CreateGoogleOauthState)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JWTAuthorization)

		// Users
		r.Get("/users/@me", userHandler.GetMe)
		r.Patch("/users/@me", userHandler.UpdateMe)

		// Agents
		r.Post("/agents/options", agentHandler.SuggestOptions)
		r.Post("/agents/chat", agentHandler.Chat)
		r.Post("/agents/{calendarId}/chat", agentHandler.CalendarChat)

		// User Profiles
		r.Put("/users/avatars", userHandler.UploadAvatar)

		// Github
		r.Post("/github/apps/state", githubHandler.CreateAppState)
		r.Post("/github/oauth/state", githubHandler.CreateOauthState)
		r.Post("/github/calendars/{calendarId}/review-requests", githubHandler.GetReviewRequests)
		r.Get("/github/calendars/{calendarId}/review-load-status", githubHandler.GetReviewLoadStatus)
		r.Get("/github/calendars/{calendarId}/milestones", githubHandler.GetMilestones)
		r.Get("/github/oauth/status", githubHandler.IsLinkedUser)
		r.Get("/github/calendars/{calendarId}/installation-status", githubHandler.IsInstalledGithubApp)

		// Calendars
		r.Post("/calendars", calendarHandler.CreateCalendar)
		r.Post("/calendars/images", calendarHandler.UploadImage)
		r.Get("/calendars/{calendarId}", calendarHandler.GetCalendar)
		r.Get("/calendars", calendarHandler.GetCalendars)
		r.Patch("/calendars/{calendarId}", calendarHandler.UpdateCalendar)
		r.Delete("/calendars/{calendarId}", calendarHandler.DeleteCalendar)

		// Members
		r.Get("/calendars/{calendarId}/members", calendarHandler.GetMembers)
		r.Post("/calendars/{calendarId}/members", calendarHandler.CreateMembers) // 使わない
		r.Post("/calendars/{calendarId}/invitations", calendarHandler.CreateInvitations)
		r.Patch("/calendars/{calendarId}/members", calendarHandler.JoinMember) // 使わない
		r.Post("/calendars/{calendarId}/join", calendarHandler.JoinMemberWithToken)
		r.Delete("/calendars/{calendarId}/members", calendarHandler.RejectMember)

		// Events
		r.Post("/calendars/{calendarId}/events", calendarHandler.CreateEvent)
		r.Patch("/calendars/{calendarId}/events/{eventId}", calendarHandler.UpdateEvent)
		r.Get("/calendars/{calendarId}/events", calendarHandler.ListGetEvents)
		r.Get("/events/todays", calendarHandler.GetToday)
		r.Get("/events/months", calendarHandler.GetByMonth)
		r.Get("/events/search", calendarHandler.SearchEvents)
		r.Delete("/calendars/{calendarId}/events/{eventId}", calendarHandler.DeleteEvent)

		// Options
		r.Get("/options", optionHandler.GetList)

		// Notices
		r.Get("/notices", noticeHttpHandler.GetNotices)

		// Scheduler
		r.Post("/calendars/{calendarId}/schedulers", schedulerHandler.CreateScheduler)
		r.Post("/calendars/{calendarId}/schedulers/{schedulerId}/attendance", schedulerHandler.AddAttendanceHandler)
		r.Get("/calendars/{calendarId}/schedulers", schedulerHandler.GetAllScheduler)
		r.Get("/calendars/{calendarId}/schedulers/{schedulerId}", schedulerHandler.GetScheduler)
		r.Get("/calendars/{calendarId}/schedulers/{schedulerId}/attendance", schedulerHandler.GetAttendance)
		r.Get("/schedulers/{schedulerId}/result", schedulerHandler.GetResult)

		// Todos
		r.Get("/calendars/{calendarId}/todos", todoHandler.GetList)
		r.Post("/calendars/{calendarId}/todos", todoHandler.CreateList)
		r.Patch("/calendars/{calendarId}/todos/{todoListId}", todoHandler.UpdateList)
		r.Post("/calendars/{calendarId}/todos/{todoListId}/items", todoHandler.AddItem)
		r.Patch("/calendars/{calendarId}/todos/{todoListId}/items/{todoItemId}", todoHandler.UpdateItem)
	})

	// Start Serving
	logrus.Info("Start serving :8000")
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
		logrus.WithField("host", host).Info("psql connection failed")
		panic(err.Error())
	}
	logrus.WithField("host", host).Info("psql connected")
	return db
}

func GetMinIOClient() *minio.Client {
	endpoint, ok := os.LookupEnv("MINIO_ENDPOINT")
	if !ok {
		panic("'MINIO_ENDPOINT' is not set'")
	}
	if os.Getenv("MINIO_FROM_IAM") == "1" {
		customTransport, err := minio.DefaultTransport(true)
		if err != nil {
			panic(err)
		}
		skipVerify := os.Getenv("MINIO_TLS_INSECURE_SKIP_VERIFY") == "1"
		customTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: skipVerify,
		}
		client, err := minio.New(endpoint, &minio.Options{
			Region:    getMinioRegion(),
			Secure:    true,
			Creds:     credentials.NewIAM(""),
			Transport: customTransport,
		})
		if err != nil {
			panic(err)
		}
		return client
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
		logrus.WithField("address", endpoint).Info("minio client connection failed")
		panic(err.Error())
	}
	logrus.WithField("address", endpoint).Info("minio client connected")
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
	opts := &redis.Options{
		Addr:     endpoint,
		Password: password,
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	}
	if os.Getenv("REDIS_TLS") == "1" {
		logrus.Info("REDIS_TLS enabled")
		skipVerify := os.Getenv("REDIS_TLS_INSECURE_SKIP_VERIFY") == "1"
		opts.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: skipVerify,
		}
	}
	client := redis.NewClient(opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logrus.WithError(err).Error("redis connection error")
		panic("redis connection failed")
	}
	logrus.WithField("address", endpoint).Info("redis connected")
	return client
}

func getEmailAddress() string {
	email, ok := os.LookupEnv("SENDER_EMAIL_ADDRESS")
	if !ok {
		panic("'SENDER_EMAIL_ADDRESS' is not set")
	}
	return email
}

func GetDialer() *mail.Dialer {
	host, ok := os.LookupEnv("EMAIL_HOST")
	if !ok {
		panic("'EMAIL_HOST' is not set")
	}
	portStr, ok := os.LookupEnv("EMAIL_PORT")
	if !ok {
		panic("'EMAIL_PORT' is not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("'EMAIL_PORT' is not int")
	}
	user, ok := os.LookupEnv("EMAIL_USER")
	if !ok {
		panic("'EMAIL_USER' is not set")
	}
	password, ok := os.LookupEnv("EMAIL_PASSWORD")
	if !ok {
		panic("'EMAIL_PASSWORD' is not set")
	}
	return mail.NewDialer(host, port, user, password)
}

func GetGenAIClient() *genai.Client {
	apiKey, ok := os.LookupEnv("AGENT_API_KEY")
	if !ok {
		panic("'AGENT_API_KEY' is not set")
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		panic(err)
	}
	return client
}

func GetOpenRouter() *openrouter.OpenRouter {
	model, ok := os.LookupEnv("AGENT_MODEL")
	if !ok {
		panic("'AGENT_MODEL' is not set")
	}
	apiKey, ok := os.LookupEnv("AGENT_API_KEY")
	if !ok {
		panic("'AGENT_API_KEY' is not set")
	}
	// Initialize LLM
	openRouter := openrouter.NewOpenRouter(apiKey)
	openRouter.SetModel(model)
	providers, ok := os.LookupEnv("AGENT_MODEL_PROVIDERS")
	if ok {
		providerList := strings.Split(providers, ",")
		openRouter.SetProviderOrder(providerList)
	}
	return openRouter
}

// 暗号キー
func getRedisEncryptionKey() []byte {
	key, ok := os.LookupEnv("REDIS_ENCRYPTION_KEY")
	if !ok {
		panic("'REDIS_ENCRYPTION_KEY' is not set")
	}
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic("failed to decode REDIS_ENCRYPTION_KEY: " + err.Error())
	}
	return decoded
}

func getInstallationIdEncryptionKey() []byte {
	key, ok := os.LookupEnv("INSTALLATION_ID_ENCRYPTION_KEY")
	if !ok {
		panic("'INSTALLATION_ID_ENCRYPTION_KEY' is not set")
	}
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic("failed to decode INSTALLATION_ID_ENCRYPTION_KEY: " + err.Error())
	}
	return decoded
}

func getGoogleOauthClientId() string {
	clientId, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID")
	if !ok {
		panic("'GOOGLE_OAUTH_CLIENT_ID' is not set")
	}
	return clientId
}

func getGoogleOauthClientSecret() string {
	clientSecret, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_SECRET")
	if !ok {
		panic("'GOOGLE_OAUTH_CLIENT_SECRET' is not set")
	}
	return clientSecret
}

func getGoogleOauthRedirectUri() string {
	redirectUri, ok := os.LookupEnv("GOOGLE_OAUTH_REDIRECT_URI")
	if !ok {
		panic("'GOOGLE_OAUTH_REDIRECT_URI' is not set")
	}
	return redirectUri
}
