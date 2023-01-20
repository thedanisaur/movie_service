package main

import (
	"log"
	"movie_service/db"
	"movie_service/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	log.Println("Starting Movie Service...")
	app := fiber.New()
	defer db.GetInstance().Close()

	// Add CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Non Authenticated routes
	app.Get("/timeline", handlers.GetTimeline)
	app.Get("/trackers", handlers.GetTrackers)
	app.Get("/ratings", handlers.GetTrackers)

	// JWT Authentication routes

	app.Listen(":1234")
}
