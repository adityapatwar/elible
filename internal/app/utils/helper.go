package utils

import (
	"crypto/rand"
	"encoding/hex"
	"mime/multipart"
	"path/filepath"
	"strings"
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
