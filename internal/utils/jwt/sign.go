package jwt

import (
	"dengovie/internal/web"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Sign(data ...any) (string, error) {
	initOnce()

	claims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	if len(data)%2 != 0 {
		return "", fmt.Errorf("invalid length of data")
	}

	var lastKey web.JWTKey
	for i := range data {
		if i%2 == 0 {
			t, ok := data[i].(web.JWTKey)
			if !ok {
				return "", fmt.Errorf("invalid type of data element at index %d (expected web.JWTKey)", i)
			}
			lastKey = t
			continue
		}
		claims[string(lastKey)] = data[i]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return tokenString, nil
}
