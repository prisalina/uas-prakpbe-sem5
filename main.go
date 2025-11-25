package main

import (
	"log"
	"os"
	"uas-pbe-praksem5/config"
	"uas-pbe-praksem5/database"
	"uas-pbe-praksem5/route"
	"uas-pbe-praksem5/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	utils.SetSecret(os.Getenv("JWT_SECRET"))
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN not set")
	}
	database.ConnectPostgres(dsn)

	app := fiber.New()
	route.RegisterAll(app, database.PG)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Listening on port", port)
	log.Fatal(app.Listen(":" + port))
}
