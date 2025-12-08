package storage

import "os"

func GetImageStorageBaseUrl() string {
	baseUrl, ok := os.LookupEnv("IMAGE_STORAGE_BASE_URL")
	if !ok {
		return "http://localhost:9000"
	}
	return baseUrl
}
