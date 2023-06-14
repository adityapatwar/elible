package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

func IsImage(file *multipart.FileHeader) bool {
	// Anda bisa memperluas atau mengubah daftar ini sesuai dengan kebutuhan Anda
	supportedExtensions := []string{".jpg", ".jpeg", ".png"}

	extension := strings.ToLower(filepath.Ext(file.Filename))

	for _, supportedExtension := range supportedExtensions {
		if extension == supportedExtension {
			return true
		}
	}

	return false
}

func RandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ParseDate(dateString string) time.Time {
	const layout = "01-02-06" // Update the layout to match the format of dateString
	t, err := time.Parse(layout, dateString)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return time.Time{} // Returning zero value of time.Time in case of error
	}
	return t
}
