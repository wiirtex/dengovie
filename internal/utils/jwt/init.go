package jwt

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
)

var jwtKey []byte

func getValueFromEnv(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		log.Fatal(fmt.Sprintf("%v environment variable not found", key))
	}

	return value
}

var initEnvsInited = atomic.Bool{}

type EnvKey string

var allEnvKeys = []string{
	"JWT_TOKEN",
}

func InitEnvs(envs map[string]string) {
	if !initEnvsInited.Load() {
		for _, key := range allEnvKeys {
			if value, in := envs[key]; !in {
				jwtKey = []byte(getValueFromEnv(key))
			} else {
				jwtKey = []byte(value)
			}
		}
	} else {
		log.Fatal(fmt.Sprintf("error: InitEnvs was called second time: %v", allEnvKeys))
	}
}
