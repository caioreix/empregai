package token

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"go-api/pkg/config"
)

// errors
var (
	ErrInvalidJWT          = errors.New("invalid jwt")
	ErrInvalidJWTSignature = errors.New("invalid jwt signature")
)

type Claims struct {
	Email string
	ID    string
	jwt.StandardClaims
}

// GenerateJWT generate a new token with claims
func GenerateJWT(email, id string, duration time.Duration, cfg *config.Config) (string, error) {
	claims := Claims{
		Email: email,
		ID:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(cfg.Server.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ExtractJWT from the request
func ExtractJWT(r *http.Request) (jwt.MapClaims, error) {
	tokenStr := extractBearerToken(r)

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (jwtKey any, err error) {
		return jwtKey, err
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrInvalidJWTSignature
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
