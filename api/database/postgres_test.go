package database

import (
	"fmt"
	"testing"
)

// DB server must be running
func TestCreatePostgresClient(t *testing.T) {
	// Local connection - correct
	dbTest, err := CreatePostgresClient("host=0.0.0.0 user=postgres password=pass dbname=shorturl port=5433 sslmode=disable")

	if dbTest == nil {
		t.Error("Postgres not connected")
	}

	if err != nil {
		fmt.Println(err)
		t.Error("Error while connecting")
	}
}