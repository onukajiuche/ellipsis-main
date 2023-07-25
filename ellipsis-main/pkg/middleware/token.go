package middleware

import (
	"brief/internal/config"
	"brief/internal/constant"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string
	Role  int
	jwt.RegisteredClaims
}

func CreateToken(id, email string, role int) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    constant.AppName,
			ID:        id,
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(config.GetConfig().SecretKey))
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Validate issuer
		if iss, err := token.Claims.GetIssuer(); iss != constant.AppName || err != nil {
			return nil, fmt.Errorf("unknown issuer: %v", token.Header["iss"])
		}

		return []byte(config.GetConfig().SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("expired token")
		} else {
			return nil, errors.New("invalid token")
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
