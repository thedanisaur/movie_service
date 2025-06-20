package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"movie_service/db"
	"movie_service/handlers"
	"movie_service/security"
	"movie_service/types"
	"movie_service/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
)

func AuthorizationMiddleware(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", util.GetFunctionName(AuthorizationMiddleware), txid.String())

		err := security.ValidateJWT(config)(c)
		if err != nil {
			log.Printf("Failed to Validate JWT\n%s\n", err)
			err_string := fmt.Sprintf("Unauthorized: %s\n", txid.String())
			return c.Status(fiber.StatusInternalServerError).SendString(err_string)
		}
		return c.Next()
	}
}

func loadConfig(config_path string) (types.Config, error) {
	var config types.Config
	config_file, err := os.Open(config_path)
	if err != nil {
		return config, err
	}
	defer config_file.Close()
	jsonParser := json.NewDecoder(config_file)
	jsonParser.Decode(&config)
	return config, nil
}

func main() {
	log.Println("Starting Movie Service...")
	config, err := loadConfig("./config.json")
	if err != nil {
		log.Printf("Error opening config, cannot continue: %s\n", err.Error())
		return
	}
	app := fiber.New()
	database, err := db.GetInstance()
	if err != nil {
		log.Printf(err.Error())
	} else {
		defer database.Close()
	}

	// Add CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.App.Cors.AllowOrigins, ","),
		AllowHeaders:     strings.Join(config.App.Cors.AllowHeaders, ","),
		AllowCredentials: config.App.Cors.AllowCredentials,
	}))

	// Add Rate Limiter
	var middleware limiter.LimiterHandler
	if config.App.Limiter.LimiterSlidingMiddleware {
		middleware = limiter.SlidingWindow{}
	} else {
		middleware = limiter.FixedWindow{}
	}
	app.Use(limiter.New(limiter.Config{
		Max:                    config.App.Limiter.Max,
		Expiration:             time.Duration(config.App.Limiter.Expiration),
		LimiterMiddleware:      middleware,
		SkipSuccessfulRequests: config.App.Limiter.SkipSuccessfulRequests,
	}))

	// Non Authenticated routes
	app.Get("/movies", handlers.GetMovies)
	app.Get("/movie_trackers/:tracker_id", handlers.GetMovieTrackersByID)
	app.Get("/movie_trackers/:movie_name/:username", handlers.GetMovieTrackers)
	app.Get("/ratings", handlers.GetRatings)
	app.Get("/series", handlers.GetSeries)
	app.Get("/timeline", handlers.GetTimeline)
	app.Get("/trackers", handlers.GetTrackers)

	// JWT Authentication routes
	app.Post("/movies/:series", AuthorizationMiddleware(config), handlers.PostMovie)
	app.Post("/movie_trackers/:username", AuthorizationMiddleware(config), handlers.PostMovieTrackers)
	app.Post("/series", AuthorizationMiddleware(config), handlers.PostSeries)
	app.Post("/trackers", AuthorizationMiddleware(config), handlers.PostTrackers)
	app.Post("/vote", AuthorizationMiddleware(config), handlers.PostVote)

	port := fmt.Sprintf(":%d", config.App.Host.Port)
	if config.App.Host.UseTLS {
		err = app.ListenTLS(port, config.App.Host.CertificatePath, config.App.Host.KeyPath)
	} else {
		log.Println("Warning - not using TLS")
		err = app.Listen(port)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}
