package main

import (
	"log"
	"os"

	"gihtub.com/ahdirmai/shorten-url-go/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err) // Menangani error godotenv dengan benar
	// }

	// println(os.Getenv("APP_PORT"))

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Port default jika APP_PORT tidak diatur
	}

	// println(port)
	log.Fatal(app.Listen(":" + port))
}
