package helpers

import (
	"strings"
)

// Generalize URL's to one format
func GeneralizeURL(url string) string {
	if url[:4] != "http" {
		genURL := strings.Replace(url, "www.", "", 1)
		trimGenURL := strings.TrimSuffix(genURL, "/")
		return "http://" + trimGenURL 
	}
	
	if url[:5] == "https" {
		genURL := strings.Replace(url, "www.", "", 1)
		trimGenURL := strings.TrimSuffix(genURL, "/")
		return strings.Replace(trimGenURL, "https", "http", 1)
	}

	genURL := strings.Replace(url, "www.", "", 1)
	trimGenURL := strings.TrimSuffix(genURL, "/")
	return trimGenURL
}