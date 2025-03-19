package middlewares

import (
	"dengovie/internal/domain"
	"dengovie/internal/utils/jwt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(ctx *gin.Context) {

	key, err := ctx.Cookie("access-token")
	if err != nil {
		if errAbort := ctx.AbortWithError(http.StatusUnauthorized, err); errAbort != nil {
			log.Printf("ctx.AbortWithError: %v\n", errAbort)
		}
	}

	jwtData, err := jwt.VerifyJWT(key)
	if err != nil {
		log.Println("verifyJWT:", err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	ctx.Set(domain.UserIDKey, jwtData)

	ctx.Next()
}
