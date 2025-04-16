package env

import (
	"fmt"
)

type Key string

const (
	KeyJwtToken Key = "JWT_TOKEN"
)

var allEnvKeys = []Key{
	KeyJwtToken,
}

var envVars = map[Key]string{}

func GetEnv(key Key) (string, error) {
	InitEnvs(nil)
	if val, ok := envVars[key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("environment variable %s not found", key)
}
