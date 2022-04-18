package main

import (
	"log"

	"devcode-pos/models"
	"devcode-pos/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	models.Conn()
}

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: true,
	})

	router.SetupRoutes(app)
	app.Listen(":3030")
}
