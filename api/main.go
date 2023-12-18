package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/lackingworth/Go-URL-Short-Ozon/routes"
)
var inMemory bool

// Setting up routes
func setupRoutes(app *fiber.App) {
	if inMemory {
		app.Get("/:url", routes.ResolveURL)
		app.Post("/api", routes.ShortenURL)
	} else {
		app.Get("/:url", routes.ResolveURLDB)
		app.Post("/api", routes.ShortenURLDB)
	}

}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	arg := os.Args[1]

	// Catch arguments from cli
	if strings.ToLower(arg) == "memory" || arg == ""{
		inMemory = true
	} else if strings.ToLower(arg) == "db" {
		inMemory = false
	}

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	app := fiber.New()

	// Init logger
	app.Use(logger.New()) 

	// Init rate limiter for api calls
	app.Use(limiter.New(limiter.Config{
		Max: 					20,
		Expiration: 			10 * time.Second,
		LimiterMiddleware: 		limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	setupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT"))) 
}