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
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Non Authenticated routes
	app.Get("/movies", handlers.GetMovies)
	app.Get("/ratings", handlers.GetRatings)
	app.Get("/series", handlers.GetSeries)
	app.Get("/timeline", handlers.GetTimeline)
	app.Get("/trackers", handlers.GetTrackers)

	// JWT Authentication routes
	app.Post("/movies/:series", handlers.PostMovie)
	app.Post("/series", handlers.PostSeries)
	app.Post("/trackers", handlers.PostTrackers)

	app.Listen(":1234")
}
