package dengovie

import (
	"dengovie/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

func getUserID(ctx *gin.Context) (domain.UserID, error) {
	userIDRaw, in := ctx.Get(domain.UserIDKey)
	if !in {
		log.Println("userID is not set")
		return 0, errors.New("userID is not set")
	}

	userID := userIDRaw.(domain.UserID)
	return userID, nil
}
