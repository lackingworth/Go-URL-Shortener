package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/lackingworth/Go-URL-Short-Ozon/database"
	"github.com/lackingworth/Go-URL-Short-Ozon/helpers"
	m "github.com/lackingworth/Go-URL-Short-Ozon/models"
)

func ShortenURL(c *fiber.Ctx) error {
	body := new(m.RequestR)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error":"Cannot parse request JSON"})
	}

	// Rate limiter - 10 times in 30 minutes from the same IP
	r2 := database.CreateClient(1)	
	defer r2.Close()
	value, err := r2.Get(database.Ctx, c.IP()).Result() // Find IP in db

	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), 10, 30*60*time.Second).Err() // Add IP to db if not found
	} else {
		value, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(value)
		
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"Error":"Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}

	}

	// URL validation
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error":"Invalid URL"})
	}

	// Domain inf loop prevention
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"Error":"Infinite loop caught"})
	}

	// Generalize links
	body.URL = helpers.GeneralizeURL(body.URL)

	
	r := database.CreateClient(0)
	defer r.Close()
	
	// Check if URL already exists in db
	shortVal, err := r.Get(database.Ctx, body.URL).Result()
	
	if err != redis.Nil {

		if body.Expiry == 0 {
			ttl, _ := r.TTL(database.Ctx, body.URL).Result()
			body.Expiry = ttl / time.Nanosecond / time.Minute
		}

		res := m.ResponseR{
			URL:					body.URL,
			CustomShort:			"",
			Expiry:					body.Expiry,
			XRateRemaining:			120,
			XRateLimitReset: 		30,
		}
		
		r2.Decr(database.Ctx, c.IP())

		value, _ = r2.Get(database.Ctx, c.IP()).Result()
		res.XRateRemaining, _ = strconv.Atoi(value)

		ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
		res.XRateLimitReset = ttl / time.Nanosecond / time.Minute
		res.CustomShort = os.Getenv("DOMAIN") + "/" + shortVal

		return c.Status(fiber.StatusOK).JSON(res)
	}

	// Custom shortened URL generation
	var id string

	if body.CustomShort == "" {
		id = helpers.GenerateRandomString(10)
	} else {
		id = body.CustomShort
	}
	value, _ = r.Get(database.Ctx, id).Result() // Checking if this short url already exists
	
	
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Error":"This URL short is already in use"})
	}
	
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, body.URL, id, body.Expiry*3600*time.Second).Err() // Add URL to db
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Unable to connect to server"})
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err() // Add short URL to db

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Unable to connect to server"})
	}

	res := m.ResponseR{
		URL:					body.URL,
		CustomShort:			"",
		Expiry:					body.Expiry,
		XRateRemaining:			120,
		XRateLimitReset: 		30,
	}

	r2.Decr(database.Ctx, c.IP())

	value, _ = r2.Get(database.Ctx, c.IP()).Result()
	res.XRateRemaining, _ = strconv.Atoi(value)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	res.XRateLimitReset = ttl / time.Nanosecond / time.Minute
	res.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(res)
}

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()
	value, err := r.Get(database.Ctx, url).Result() // Getting URL from Redis

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error":"Short URL not found in the memory"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Cannot connect to Redis"})
	}

	value = helpers.GeneralizeURL(value)

	// Increments counter
	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")
	
	// To get the url in response instead of redirect use 
	// return c.Status(fiber.StatusOK).JSON(value)
	return c.Redirect(value, 301) 
}