package dengovie

import (
	"dengovie/internal/domain"
	"dengovie/internal/utils/jwt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListUserGroups godoc
//
//	@Summary      Получение информации о профиле пользователя
//	@Description  get-me
//	@Accept       json
//	@Produce      json
//	@Success		200		{object}	string			"ok"
//	@Router        /user [get]
func (c *Controller) GetMe(ctx *gin.Context) {
	usr, exist := ctx.Get(domain.UserIDKey)
	if !exist {
		log.Println("no usr data in ctx, check that route contains CheckAuth middleware")
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	jwtData, ok := usr.(jwt.JWTData)
	if !ok {
		log.Println("usr is not of type jwtData")
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, jwtData)
}
