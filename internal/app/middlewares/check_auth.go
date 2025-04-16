package middlewares

import (
	"dengovie/internal/domain"
	"dengovie/internal/utils/jwt"
	"dengovie/internal/web"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(ctx *gin.Context) {
	key, err := ctx.Cookie("access-token")
	if err != nil {
		if errAbort := ctx.AbortWithError(http.StatusUnauthorized, err); errAbort != nil {
			log.Printf("ctx.AbortWithError: %v\n", errAbort)
			return
		}
	}

	jwtData, err := jwt.VerifyJWT(key)
	if err != nil {
		log.Println("verifyJWT:", err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID, ok := jwtData[web.JWTUserIDKey].(domain.UserID)
	if !ok {
		log.Println("userID is not int64")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set(domain.UserIDKey, userID)

	ctx.Next()
}
