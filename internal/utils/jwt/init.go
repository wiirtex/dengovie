package jwt

import (
	"log"
	"os"
)

var jwtKey []byte

func init() {
	token, found := os.LookupEnv("JWT_TOKEN")
	if !found {
		log.Fatal("JWT_TOKEN environment variable not found")
	}

	jwtKey = []byte(token)
}
