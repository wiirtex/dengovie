package jwt

import (
	"dengovie/internal/utils/env"
	"fmt"
)

var jwtKey []byte

func init() {
	value, err := env.GetEnv(env.KeyJwtToken)
	if err != nil {
		panic(fmt.Errorf("JWT token env error: %w", err))
	}

	jwtKey = []byte(value)
}
