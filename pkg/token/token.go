package token

import (
	"strings"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	errors "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID  uuid.UUID  `json:"user_id"`
	Role    string     `json:"role"`
	GroupID *uuid.UUID `json:"group_id,omitempty"`
	jwt.RegisteredClaims
}

type ClaimsRefresh struct {
	UserID   uuid.UUID `json:"user_id"`
	UserType string    `json:"user_type"`
	jwt.RegisteredClaims
}

type Token interface {
	ParseToken(token string) (*Claims, error)
	ExtractToken(c *gin.Context) string
	GenerateToken(userID uuid.UUID, group_id *uuid.UUID, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, user_type string) (*model.Refresh, error)
}

type jwtService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttl time.Duration) Token {
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

func (s *jwtService) GenerateToken(userID uuid.UUID, group_id *uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID:  userID,
		Role:    role,
		GroupID: group_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(timezone.Dushanbe).Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(timezone.Dushanbe)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *jwtService) GenerateRefreshToken(userID uuid.UUID, user_type string) (*model.Refresh, error) {
	claims := ClaimsRefresh{
		UserID:   userID,
		UserType: user_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(timezone.Dushanbe).Add(24 * time.Hour * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(timezone.Dushanbe)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refresh, err := token.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	r := &model.Refresh{
		UserID:    userID,
		UserType:  user_type,
		Token:     refresh,
		ExpiresAt: time.Now().In(timezone.Dushanbe).Add(time.Hour * 24 * 30),
	}
	return r, nil
}
