package database

import (
	"fmt"

	m "github.com/lackingworth/Go-URL-Short-Ozon/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Dsn string = "host=postgres user=postgres password=pass dbname=shorturl port=5432 sslmode=disable"

func CreatePostgresClient(dsn string) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		fmt.Println(err)
		return db, err
	}

	err = db.AutoMigrate(&m.ResponseP{})

	if err != nil {
		fmt.Println(err)
		return db, err
	}
	
	return db, nil
}