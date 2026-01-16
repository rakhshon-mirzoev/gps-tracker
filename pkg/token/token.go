package token

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type Token interface {
	ParseToken(token string) (*Claims, error)
	ExtractToken(c *gin.Context) string
	GenerateToken(entityID uuid.UUID, entityRole string) (string, error)
}

type jwtService struct {
	secret []byte
	ttl    time.Duration
}

func InitTokenService(secret string, ttl time.Duration) Token {
	return &jwtService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *jwtService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrorSignMethod
		}
		return s.secret, nil
	},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrorInvalidToken
	}

	return claims, nil
}

func (s *jwtService) ExtractToken(c *gin.Context) string {
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

func (s *jwtService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}
