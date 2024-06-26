package crypto

import (
	"fmt"
	"time"

	configs "halo_suster/cfg"
	"halo_suster/db/entities"

	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     string
	jwtSecretOnce sync.Once
)

// initJWTSecret initializes the JWT secret once to avoid repeated environment variable lookups.
func initJWTSecret() {
	jwtSecret = configs.GetEnvOrDefault("JWT_SECRET", "sec")
}

// GenerateToken generates a JWT token with the provided user details.
func GenerateToken(id, nip, name string) (string, error) {
	jwtSecretOnce.Do(initJWTSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.JWTClaims{
		Id:   id,
		Nip:  nip,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken verifies a JWT token and returns the payload if valid.
func VerifyToken(tokenString string) (*entities.JWTPayload, error) {
	jwtSecretOnce.Do(initJWTSecret)

	claims := &entities.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	payload := &entities.JWTPayload{
		Id:   claims.Id,
		Nip:  claims.Nip,
		Name: claims.Name,
	}

	return payload, nil
}
