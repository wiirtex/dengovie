package env

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var initOnce sync.Once

func getValueFromEnv(key Key) string {
	value, found := os.LookupEnv(string(key))
	if !found {
		log.Fatal(fmt.Sprintf("%v environment variable not found", key))
	}

	return value
}

func InitEnvs(envs map[Key]string) {
	initOnce.Do(func() {
		for _, key := range allEnvKeys {
			if value, in := envs[key]; !in {
				envVars[key] = getValueFromEnv(key)
			} else {
				envVars[key] = value
			}
		}
	})
}
