package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

func main() {
	_ = godotenv.Load()

	database := getPostgresDB()
	migrations := migrate.FileMigrationSource{
		Dir: "db/migrate",
	}
	n, err := migrate.Exec(database.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Applied %d migrations.\n", n)
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
