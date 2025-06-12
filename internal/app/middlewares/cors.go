package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

var allowList = map[string]bool{
	"http://localhost:5173":  true, // local development (frontend)
	"http://backend.ingress": true, // production (ingress for backend)
	"http://localhost:8080":  true, // local development (swagger to ingress)
}

func CORSMiddleware(c *gin.Context) {
	log.Println("origin:", c.Request.Header.Get("Origin"))
	if origin := c.Request.Header.Get("Origin"); allowList[origin] {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	}
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204) // No Content
		return
	}

	c.Next()
}
