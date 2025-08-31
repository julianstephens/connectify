package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx stdlib driver
)

func NewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	// Open as database/sql but backed by pgx driver
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pooling as needed
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Quick ping with timeout
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctxPing); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
