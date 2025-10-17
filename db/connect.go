package db

import (
	"database/sql"
	"fmt"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	_ "github.com/lib/pq"
)

func Connect(dbName, dbUser, dbPass, dbHost string, dbPort int) *sql.DB {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		logger.Log.Fatalf("Unable to connec to the database: %v", err)
	}

	return db
}
