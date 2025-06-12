package dengovie

import (
	"context"
	storeTypes "dengovie/internal/store/types"
	"dengovie/internal/utils/jwt"
	"dengovie/internal/web"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestCode godoc
//
//	@Summary      Запросить код для входа
//	@Description  request-otp-code
//	@Accept       json
//	@Produce      json
//	@Param			telegramAlias 	body 	string	true	 "telegramAlias"
//	@Success		200		{string}	string			"ok, код отправлен"
//	@Failure		400		{object}	web.APIError	"данные невалидные"
//	@Failure		404		{object}	web.APIError	"клиент не зарегистрирован в боте"
//	@Router        /auth/request_code [post]
func (c *Controller) RequestCode(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// Login godoc
//
//	@Summary      Зайти с telegramAlias и коду. На время тестирования подходит 111
//	@Description  login-with-otp-code
//	@Accept       json
//	@Produce      json
//	@Param			body 	body 	LoginRequest	true	 "body"
//	@Success		200		{string}	string			"ok"
//	@Failure		400		{object}	web.APIError	"данные невалидные"
//	@Router        /auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {

	req := LoginRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		log.Println("error ShouldBindBodyWithJSON:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	if req.OTP != "111" {
		log.Println("invalid otp")
		ctx.JSON(http.StatusForbidden, web.APIError{
			ErrorReason: web.InvalidOTP,
		})
		return
	}

	user, err := c.storage.GetUserByAlias(context.TODO(), storeTypes.GetUserByAliasInput{
		Alias: req.Alias,
	})
	if err != nil {
		log.Println("storage.GetUserByAlias:", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	signetJWT, err := jwt.Sign(web.JWTUserIDKey, user.ID)
	if err != nil {
		log.Println("jwt.Sign:", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("access-token", signetJWT, 0, "/", "dengovie.ingress", false, true)
}

// Logout godoc
//
//	@Summary      Выйти из профиля
//	@Description  logout
//	@Success		200		{string}	string			"ok"
//	@Router        /auth/logout [post]
func (c *Controller) Logout(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("access-token", "", 0, "/", "dengovie.ingress", false, true)
}

type LoginRequest struct {
	Alias string `json:"alias"`
	OTP   string `json:"otp"`
}
