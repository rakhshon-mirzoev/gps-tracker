package utils

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"golang.org/x/crypto/bcrypt"
)

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

func SaveBase64Image(data string, name *string, entity string, id uuid.UUID) (string, error) {
	if name == nil {
		return "", errs.ErrorNameIsNil
	}

	if strings.TrimSpace(data) == "" {
		return "", errs.ErrorEmptyBase64
	}

	mimePart := strings.SplitN(data, ";", 2)[0]
	if !strings.HasPrefix(mimePart, "data:image/") {
		return "", fmt.Errorf("invalid MIME type: %s", mimePart)
	}
	fileType := strings.TrimPrefix(mimePart, "data:image/")

	commaIdx := strings.Index(data, ",")
	if commaIdx == -1 {
		return "", fmt.Errorf("invalid base64 image format")
	}
	base64Data := data[commaIdx+1:]

	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	safeName := strings.ReplaceAll(*name, " ", "_")
	filename := fmt.Sprintf("%s_%d.%s", safeName, time.Now().Unix(), fileType)

	fullPath, err := PathCreator(filename, entity, *name, id)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, decoded, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}

func PathCreator(fileName, entity, dir string, id uuid.UUID) (string, error) {
	uploadDir := "uploads"
	dir = fmt.Sprintf("%s/%s/%s/%d", uploadDir, entity, dir, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	savePath := filepath.Join(dir, fileName)

	return savePath, nil
}

func ValidateFileType(file *multipart.FileHeader) error {
	allowedExtensions := []string{".png", ".jpg", ".jpeg"}
	ext := strings.ToLower(filepath.Ext(file.Filename))

	for _, t := range allowedExtensions {
		if ext == t {
			return nil
		}
	}
	return errs.ErrorInvalidFileFormat
}
