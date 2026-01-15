package token

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type Token interface {
	ParseClaims(c *gin.Context) (uuid.UUID, string, error)
	ExtractToken(c *gin.Context) string
	GenerateToken(entityID uuid.UUID, entityRole string) (string, error)
}

type TokenService struct {
}

func InitTokenService() Token {
	return &TokenService{}
}

func (s *TokenService) ParseClaims(c *gin.Context) (uuid.UUID, string, error) {
	tokenString := s.ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		idFloat, ok := claims["id"].(uuid.UUID)
		if !ok {
			return uuid.Nil, "", fmt.Errorf("invalid id format")
		}
		role, ok := claims["role"].(string)
		if !ok {
			return uuid.Nil, "", fmt.Errorf("invalid role format")
		}

		return idFloat, role, nil
	}
	return uuid.Nil, "", fmt.Errorf("invalid token")
}

func (s *TokenService) ExtractToken(c *gin.Context) string {
	if token := c.Query("token"); token != "" {
		return token
	}
	bearerToken := c.GetHeader("Authorization")
	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func (s *TokenService) GenerateToken(entityID uuid.UUID, entityRole string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"id":         entityID,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
		"status":     "OK",
		"role":       entityRole,
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", errors.ErrorNoSecretKey
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
