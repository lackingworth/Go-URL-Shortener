package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/lackingworth/Go-URL-Short-Ozon/database"
	"github.com/lackingworth/Go-URL-Short-Ozon/helpers"
	m "github.com/lackingworth/Go-URL-Short-Ozon/models"
)

var Dsn = "host=0.0.0.0 user=postgres password=pass dbname=shorturl port=5433 sslmode=disable"

func TestShortenURLDB(t *testing.T) {
	err := godotenv.Load("../.env")

	if err != nil {
		t.Errorf("Error loading .env file")
	}

	firstReq := m.RequestP{
		URL: "www.youtube.com",
	} 
	secReq := m.RequestP{
		URL: "gibberish string",
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

		var res m.ResponseP = m.ResponseP{}
		db, _ := database.CreatePostgresClient(Dsn)
		body := new(m.RequestP)

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error":"Cannot parse request JSON"})
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
		
		// Check if URL already exists in db
		tx := db.Where("url = ?", body.URL).First(&res)

		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusOK).JSON(res) // Return found entry
		} else if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while retrieving query for url"}) // Catch unexpected error
		}

		// Custom shortened URL generation
		var id string

		// Check if user presented short url is unique
		if body.CustomShort == "api" || body.CustomShort == "api/" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"Error":"Conflict with endpoint caught"})
		}

		if body.CustomShort != "" {
			tx := db.Where("short_url = ?", body.CustomShort).First(&res)

			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				id = body.CustomShort
				urlEntry := m.ResponseP{URL: body.URL, ShortURL: os.Getenv("DOMAIN") + "/" + id} 
				tx := db.Create(&urlEntry) // Add URL to db

				if tx.Error != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while creating url query"}) 
				}
				
				return c.Status(fiber.StatusOK).JSON(urlEntry)

			} else if tx.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while retrieving query for short url"})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error":"Custom URL already in use"})
		}

		id = helpers.GenerateRandomString(10)
		urlEntry := m.ResponseP{URL: body.URL, ShortURL: os.Getenv("DOMAIN") + "/" + id} 
		tx = db.Create(&urlEntry) // Add URL to db
		
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while creating url query"}) 
		}

		return c.Status(fiber.StatusOK).JSON(urlEntry)
	})

	
	
	for _, test := range tests {
		
		req := httptest.NewRequest("POST", test.route, test.bodyT)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.desc)
		
		} 
}

func TestResolveURLDB(t *testing.T) {
	err := godotenv.Load("../.env")
	
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	
	tests := []struct {
		desc			string
		route			string
		expectedCode	int
	}{
		{
			desc: 			"GET http status 301 success",
			route: 			"R501T7OLLz", // Created beforehand
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

		var str m.ResponseP = m.ResponseP{}
		db, _ := database.CreatePostgresClient(Dsn)
		url := tests[0].route
		tx := db.Where("short_url = ?", os.Getenv("DOMAIN") + "/" + url).First(&str) // Getting URL from db
		
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error":"Short URL not found in db"})
		} else if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while searching for original url"})
		}

		str.URL = helpers.GeneralizeURL(str.URL)

		// To get the url in response instead of redirect use 
		// return c.Status(fiber.StatusOK).JSON(str.URL)
		return c.Redirect(str.URL, 301)
	})

	req1 := httptest.NewRequest("GET", "/R501T7OLLz", nil)
	resp1, _ := app.Test(req1, -1)

	assert.Equalf(t, tests[0].expectedCode, resp1.StatusCode, tests[0].desc)

	// Second suite
	app.Get(tests[1].route, func(c *fiber.Ctx) error {

		var str m.ResponseP = m.ResponseP{}
		db, _ := database.CreatePostgresClient(Dsn)
		url := tests[1].route
		tx := db.Where("short_url = ?", os.Getenv("DOMAIN") + "/" + url).First(&str) // Getting URL from db
		
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error":"Short URL not found in db"})
		} else if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while searching for original url"})
		}

		str.URL = helpers.GeneralizeURL(str.URL)

		// To get the url in response instead of redirect use 
		// return c.Status(fiber.StatusOK).JSON(str.URL)
		return c.Redirect(str.URL, 301)
	})

	req2 := httptest.NewRequest("GET", "/gibberishshortlink", nil)
	resp2, _ := app.Test(req2, -1)

	assert.Equalf(t, tests[1].expectedCode, resp2.StatusCode, tests[1].desc)
}
