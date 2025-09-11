package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func PanicCatcher(ctx *gin.Context) {

	defer func() {
		if msg := recover(); msg != nil {
			log.Printf("panic recovered: %v", msg)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}()

	ctx.Next()
}
