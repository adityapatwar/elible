package utils

import (
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
