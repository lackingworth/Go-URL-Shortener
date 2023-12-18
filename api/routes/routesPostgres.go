package routes

import (
	"errors"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/lackingworth/Go-URL-Short-Ozon/database"
	"github.com/lackingworth/Go-URL-Short-Ozon/helpers"
	m "github.com/lackingworth/Go-URL-Short-Ozon/models"
	"gorm.io/gorm"
)

// Shorten provided url - POST req
func ShortenURLDB(c *fiber.Ctx) error {
	var res m.ResponseP = m.ResponseP{}
	db, conErr := database.CreatePostgresClient(database.Dsn)

	if conErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":"Error while connecting to database"})
	}

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
}

// Redirect from short - GET req
func ResolveURLDB(c *fiber.Ctx) error {
	var str m.ResponseP = m.ResponseP{}
	db, _ := database.CreatePostgresClient(database.Dsn)
	url := c.Params("url")
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
}