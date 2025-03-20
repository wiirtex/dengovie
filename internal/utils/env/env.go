package env

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
)

var initEnvsInited = atomic.Bool{}

func getValueFromEnv(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		log.Fatal(fmt.Sprintf("%v environment variable not found", key))
	}

	return value
}

func InitEnvs(envs map[string]string) {
	if !initEnvsInited.Load() {
		for _, key := range allEnvKeys {
			value, in := envs[key]

			envValues[key] = value
			if !in {
				envValues[key] = getValueFromEnv(key)
			}
		}
	} else {
		log.Fatal(fmt.Sprintf("error: InitEnvs was called second time: %v", allEnvKeys))
	}
}
