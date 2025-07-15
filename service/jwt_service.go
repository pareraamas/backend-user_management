package service

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Service untuk generate dan verifikasi JWT (token session)
// Komentar dalam Bahasa Indonesia

type JWTService struct {
	secretKey string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secretKey: secret}
}

func (j *JWTService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
}
