package authModel

import "github.com/golang-jwt/jwt/v5"

type UserLogin struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
