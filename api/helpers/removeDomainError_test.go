package helpers

import (
	"testing"

	"github.com/joho/godotenv"
)

type addTest struct {
	incomingString 	string
	resultBool 		bool
}

var addTests = []addTest{
	{"http://www.youtube.com", true},
	{"http://localhost:3000", false},
	{"https://localhost:3000", false},
	{"www.localhost:3000", false},
	{"localhost:3000", false},
}

func TestRemoveDomainError(t *testing.T) {
	err := godotenv.Load("../.env")

	if err != nil {
		t.Errorf("Error loading .env file")
	}
	
	for _, test := range addTests {
		if output := RemoveDomainError(test.incomingString); output != test.resultBool {
			t.Errorf("Output %v not equal to expected %v", output, test.resultBool)
		}
	}
}