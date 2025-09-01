package utils

import (
	"errors"
	"fmt"
	"runs-system-user-go/config"
	authModel "runs-system-user-go/module/auth/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWT_SECRET = config.Config("JWT_SECRET")

func GenerateToken(email string, userId uuid.UUID) (string, error) {
	claims := authModel.Claims{
		UserID: uint(userId[0]),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(email),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // token valid 24h
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_SECRET))
}

func ParseToken(tokenString string) (*authModel.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authModel.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*authModel.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
