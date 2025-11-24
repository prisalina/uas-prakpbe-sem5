package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var PG *sql.DB

func ConnectPostgres(dsn string) {
	var err error
	PG, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("postgres open error: %v", err)
	}
	PG.SetMaxOpenConns(25)
	PG.SetMaxIdleConns(25)
	PG.SetConnMaxLifetime(5 * time.Minute)

	if err = PG.Ping(); err != nil {
		log.Fatalf("postgres ping error: %v", err)
	}
	log.Println("Connected to Postgres")
}
