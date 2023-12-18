package routes

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strconv"
	"time"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/go-redis/redis/v8"

	"github.com/lackingworth/Go-URL-Short-Ozon/database"
	"github.com/lackingworth/Go-URL-Short-Ozon/helpers"
	m "github.com/lackingworth/Go-URL-Short-Ozon/models"
)

var Address = "127.0.0.1:6379"

func TestShortenURL(t *testing.T) {
	err := godotenv.Load("../.env")

	if err != nil {
		t.Errorf("Error loading .env file")
	}

	firstReq := m.RequestP{
		URL: "www.generalmotors.com",
	} 
	secReq := m.RequestP{
		URL: "not a url",
	}
	thirdReq := m.RequestP{
		URL: "localhost:3000",
	}
	fourthReq := m.RequestP{
		URL: "www.vk.com",
	}
	
	marshalled1, _ := json.Marshal(firstReq)
	marshalled2, _ := json.Marshal(secReq)
	marshalled3, _ := json.Marshal(thirdReq)
	marshalled4, _ := json.Marshal(fourthReq)
	
	
	tests := []struct {
		desc			string
		route			string
		bodyT			*bytes.Reader
		expectedCode	int
	}{
		{
			desc: 			"POST http status 200 success",
			route: 			"/api",
			bodyT:			bytes.NewReader(marshalled1),
			expectedCode: 	200,
		},
		{
			desc: 			"POST http status 400 not a url",
			route: 			"/api",
			bodyT:			bytes.NewReader(marshalled2),
			expectedCode: 	400,
		},
		{
			desc: 			"POST http status 503 inf loop",
			route: 			"/api",
			bodyT:			bytes.NewReader(marshalled3),
			expectedCode: 	503,
		},
		{
			desc: 			"POST http status 404 no endpoint",
			route: 			"/testendpoint",
			bodyT:			bytes.NewReader(marshalled4),
			expectedCode: 	404,
		}, 
	}

	app := fiber.New()
	app.Post("/api", func(c *fiber.Ctx) error {
		body := new(m.RequestR)

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error":"Cannot parse request JSON"})
		}

		// Rate limiter - 10 times in 30 minutes from the same IP
		r2 := database.CreateClient(1, Address)	
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

		
		r := database.CreateClient(0, Address)
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
				XRateRemaining:			10,
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
			XRateRemaining:			10,
			XRateLimitReset: 		30,
		}

		r2.Decr(database.Ctx, c.IP())

		value, _ = r2.Get(database.Ctx, c.IP()).Result()
		res.XRateRemaining, _ = strconv.Atoi(value)

		ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
		res.XRateLimitReset = ttl / time.Nanosecond / time.Minute
		res.CustomShort = os.Getenv("DOMAIN") + "/" + id

		return c.Status(fiber.StatusOK).JSON(res)
	})

	for _, test := range tests {
		
		req := httptest.NewRequest("POST", test.route, test.bodyT)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.desc)
		
	}
}

func TestResolveURL(t *testing.T) {
	tests := []struct {
		desc			string
		route			string
		expectedCode	int
	}{
		{
			desc: 			"GET http status 301 success",
			route: 			"nSdfi347kd", // Created beforehand
			expectedCode: 	301,
		},
		{
			desc: 			"GET http status 404 url not found / no endpoint",
			route: 			"/gibberishshortlink",
			expectedCode: 	404,
		},
	}

	app := fiber.New()

	// First suite
	app.Get(tests[0].route, func(c *fiber.Ctx) error {

		url := tests[0].route
		r := database.CreateClient(0, Address)
		defer r.Close()
		value, err := r.Get(database.Ctx, url).Result() // Getting URL from Redis

		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error":"Short URL not found in the memory"})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Cannot connect to Redis"})
		}

		value = helpers.GeneralizeURL(value)

		// Increments counter
		rInr := database.CreateClient(1, Address)
		defer rInr.Close()
		_ = rInr.Incr(database.Ctx, "counter")
		
		// To get the url in response instead of redirect use 
		// return c.Status(fiber.StatusOK).JSON(value)
		return c.Redirect(value, 301) 
	})

	req1 := httptest.NewRequest("GET", "/nSdfi347kd", nil)
	resp1, _ := app.Test(req1, -1)

	assert.Equalf(t, tests[0].expectedCode, resp1.StatusCode, tests[0].desc)

	// Second suite
	app.Get(tests[1].route, func(c *fiber.Ctx) error {

		url := tests[1].route
		r := database.CreateClient(0, Address)
		defer r.Close()
		value, err := r.Get(database.Ctx, url).Result() // Getting URL from Redis

		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error":"Short URL not found in the memory"})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Cannot connect to Redis"})
		}

		value = helpers.GeneralizeURL(value)

		// Increments counter
		rInr := database.CreateClient(1, Address)
		defer rInr.Close()
		_ = rInr.Incr(database.Ctx, "counter")
		
		// To get the url in response instead of redirect use 
		// return c.Status(fiber.StatusOK).JSON(value)
		return c.Redirect(value, 301) 
	})

	req2 := httptest.NewRequest("GET", "/gibberishshortlink", nil)
	resp2, _ := app.Test(req2, -1)

	assert.Equalf(t, tests[1].expectedCode, resp2.StatusCode, tests[1].desc)
}