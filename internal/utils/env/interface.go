package env

import (
	"errors"
	"fmt"
)

var allEnvKeys = []string{
	"JWT_TOKEN",
}

var envValues = make(map[string]string)

func GetEnv(key string) (string, error) {
	if !initEnvsInited.Load() {
		return "", errors.New("environment variables not initialised")
	}
	if val, ok := envValues[key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("environment variable %s not found", key)
}
