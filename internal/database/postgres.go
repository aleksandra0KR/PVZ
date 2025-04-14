package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Postgres struct {
	db                 *sqlx.DB
	MaxIdleConnections int
	MaxOpenConnections int
}

func InitializeDBPostgres(maxIdleConnections, maxOpenConnections int) *Postgres {
	postgresDB := Postgres{
		MaxIdleConnections: maxIdleConnections,
		MaxOpenConnections: maxOpenConnections,
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable`, dbHost, dbUser, dbPassword, dbName, dbPort)
	log.Infof(connStr)

	var db *sqlx.DB
	var err error
	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", connStr)
		if err == nil {
			break
		}
		log.Warnf("failed to connect to database: %v. Retrying in %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}
	if err != nil || db == nil {
		log.Fatalf("failed to connect to database after %d retries: %v", maxRetries, err)
	}

	db.SetMaxIdleConns(postgresDB.MaxIdleConnections)
	db.SetMaxOpenConns(postgresDB.MaxOpenConnections)
	postgresDB.db = db
	log.Info("connected to Postgres DB")

	postgresDB.migrate()
	return &postgresDB
}

func (postgresDB *Postgres) migrate() {
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	pathToMigrations := os.Getenv("PATH_TO_MIGRATIONS")
	if pathToMigrations == "" {
		pathToMigrations = "./migrations"
	}
	if err := goose.Up(postgresDB.db.DB, pathToMigrations); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}

func (postgresDB *Postgres) GetDB() *sqlx.DB {
	return postgresDB.db
}
