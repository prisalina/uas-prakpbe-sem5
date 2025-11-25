package main

import (
	"log"
	"os"
	"uas-pbe-praksem5/config"
	"uas-pbe-praksem5/database"
	"uas-pbe-praksem5/utils"
)

func main() {
	config.LoadEnv()
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN not set")
	}
	database.ConnectPostgres(dsn)

	db := database.PG
	var roleID string
	err := db.QueryRow(`SELECT id FROM roles WHERE name='Admin'`).Scan(&roleID)
	if err != nil {
		log.Fatalf("failed to get admin role id: %v", err)
	}

	username := "admin"
	email := "admin@local"
	password := "123456"
	fullName := "Admin System"

	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("hash error: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users (username,email,password_hash,full_name,role_id,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5,NOW(),NOW()) ON CONFLICT (username) DO NOTHING`, username, email, hash, fullName, roleID)
	if err != nil {
		log.Fatalf("insert admin error: %v", err)
	}
	log.Println("Admin seeded. username:", username, "password:", password)
}
