package jwt

import (
	"dengovie/internal/utils/env"
	"dengovie/internal/utils/jwt/types"
	"fmt"
	"sync"
)

type jwtProcessor struct {
}

func New() types.Processor {
	return &jwtProcessor{}
}

var jwtKey []byte

var initFunc sync.Once

func initOnce() {
	initFunc.Do(func() {
		value, err := env.GetEnv(env.KeyJwtToken)
		if err != nil {
			panic(fmt.Errorf("JWT token env error: %w", err))
		}

		jwtKey = []byte(value)
	})
}
