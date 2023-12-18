package helpers

import "testing"

func TestGenerateRandomString(t *testing.T) {
	result_1 := GenerateRandomString(10)
	result_2 := GenerateRandomString(10)
	result_3 := GenerateRandomString(10)

	if result_1 == result_2 || result_1 == result_3 {
		t.Errorf("Generated strings are not unique")
	}

	if len(result_1) != 10 || len(result_2) != 10 || len(result_3) != 10 {
		t.Errorf("Incorrect length")
	}
}