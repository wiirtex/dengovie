package middlewares

import (
	"dengovie/internal/utils/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	ctx.Set("user_id", jwtData.UserID)

	ctx.Next()
}
