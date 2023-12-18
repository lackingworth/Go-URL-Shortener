package database

import (
	"fmt"

	m "github.com/lackingworth/Go-URL-Short-Ozon/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbIsConnected bool = false

func CreatePostgresClient() *gorm.DB {
	dsn := "host=postgres user=postgres password=pass dbname=shorturl port=5432 sslmode=disable" // For local testing change info to your local instance
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&m.ResponseP{})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected to Postgres")
	
	return db
}