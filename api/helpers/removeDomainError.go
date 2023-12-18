package helpers

import "os"

// Prevent infinite callback abuse
func RemoveDomainError(url string) bool {
	genURL := GeneralizeURL(url)
	
	if genURL == "http://" + os.Getenv("DOMAIN") {
		return false
	}

	return true
}