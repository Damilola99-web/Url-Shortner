package helpers

import "os"

func EnforceHTTP(url string) string {
	if url[:5] != "https" {
		return "https://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	return url != os.Getenv("DOMAIN")

}
