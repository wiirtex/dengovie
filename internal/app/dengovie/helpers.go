package dengovie

import (
	"dengovie/internal/domain"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

func getUserID(ctx *gin.Context) (domain.UserID, error) {
	userIDRaw, in := ctx.Get(domain.UserIDKey)
	if !in {
		log.Println("userID is not set")
		return 0, errors.New("userID is not set")
	}

	userID, ok := userIDRaw.(domain.UserID)
	if !ok {
		return 0, errors.New("userID is not a domain.UserID")
	}

	return userID, nil
}
