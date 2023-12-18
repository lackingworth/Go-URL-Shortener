package helpers

import (
	"math/rand"
)

// Generate random string of 10 symbols from given charset
func GenerateRandomString(length int) string {
	const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	res := make([]byte, length)
	for i := range res {
		res[i] = CHARSET[rand.Intn(len(CHARSET))]
	}
	return string(res)
}