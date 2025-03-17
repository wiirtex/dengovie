package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func Sign(data ...any) (string, error) {
	claims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	if len(data)%2 != 0 {
		return "", fmt.Errorf("invalid length of data")
	}

	var lastKey string
	for i := range data {
		if i%2 == 0 {
			t, ok := data[i].(string)
			if !ok {
				return "", fmt.Errorf("invalid type of data element at index %d (expected string)", i)
			}
			lastKey = t
			continue
		}
		claims[lastKey] = data[i]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return tokenString, nil
}
