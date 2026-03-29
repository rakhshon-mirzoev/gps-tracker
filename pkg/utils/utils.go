package utils

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const earthRadius = 6371000

const MaxFileSize = 10 * 1024 * 1024

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}

func Join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}

func TotalPage(total, pageSize int64) int64 {
	return (total + pageSize - 1) / pageSize
}

func PathCreator(fileName, dir, entity string, id uuid.UUID) (string, error) {
	uploadDir := "uploads"
	dir = fmt.Sprintf("%s/%s/%s/%d", uploadDir, dir, entity, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	savePath := filepath.Join(dir, fileName)

	return savePath, nil
}

func GetFileExtension(filename string) string {
	if idx := strings.LastIndexByte(filename, '.'); idx != -1 {
		return strings.ToLower(filename[idx:])
	}
	return ""
}

func ValidateFileSize(file *multipart.FileHeader) error {
	if file == nil {
		return errs.ErrorInvalidFileFormat
	}
	if file.Size > MaxFileSize {
		return errs.ErrorExceedsFileSize
	}
	if file.Size <= 0 {
		return errs.ErrorFileEmpty
	}
	return nil
}

func ValidateFileType(file *multipart.FileHeader) error {
	if file == nil {
		return errs.ErrorInvalidFileFormat
	}
	allowedExtensions := []string{".png", ".jpg", ".jpeg"}
	ext := strings.ToLower(filepath.Ext(file.Filename))

	for _, t := range allowedExtensions {
		if ext == t {
			return nil
		}
	}
	return errs.ErrorInvalidFileFormat
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func ValidateStoreStatus(status string) error {
	if status == "" {
		return nil
	}
	m := map[string]struct{}{
		"pending":     {},
		"approved":    {},
		"deactivated": {},
		"frozen":      {},
		"rejected":    {},
	}

	_, ok := m[status]
	if !ok {
		return errs.ErrorInvalidStatus
	}

	return nil

}

func ValidateAllPhotos(interior, exterior []string) error {
	seen := make(map[string]struct{})

	check := func(keys []string) error {
		for _, key := range keys {
			if !strings.HasPrefix(key, "stores/tmp/") {
				return fmt.Errorf("invalid photo key: %s", key)
			}

			if _, exists := seen[key]; exists {
				return fmt.Errorf("%w: %s", errs.ErrorDuplicatePhotoKey, key)
			}

			seen[key] = struct{}{}
		}
		return nil
	}

	if err := check(interior); err != nil {
		return err
	}

	if err := check(exterior); err != nil {
		return err
	}

	return nil
}
