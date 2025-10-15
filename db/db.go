package db

import (
	"database/sql"
	"fmt"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/logger"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(dbName, dbUser, dbPassword, dbHost string, dbPort int) *sql.DB {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to postgresql: %v", err)
	}

	logger.Log.Info("Successfully connected to the database")
	return conn
}
