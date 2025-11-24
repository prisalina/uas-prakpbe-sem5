package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func SetSecret(s string) {
	jwtSecret = []byte(s)
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username, roleName string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
