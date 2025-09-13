package dengovie

import (
	"dengovie/internal/domain"
	usersTypes "dengovie/internal/service/users/types"
	storeTypes "dengovie/internal/store/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMe godoc
//
//	@Summary      Получение информации о профиле пользователя
//	@Description  get-me
//	@Accept       json
//	@Produce      json
//	@Success		200		{object}	GetMeResponse			"ok"
//	@Router        /user [get]
func (c *Controller) GetMe(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := c.storage.GetUserByID(ctx, storeTypes.GetUserByIDInput{
		UserID: domain.UserID(userID),
	})
	if err != nil {
		log.Println("storage.GetUserByID err:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, GetMeResponse{
		UserID: int64(user.ID),
		Name:   user.Name,
		Alias:  user.Alias,
	})
}

type GetMeResponse struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Alias  string `json:"alias"`
}

// UpdateName godoc
//
//	@Summary      Обновление имени юзера
//	@Description  update-my-name
//	@Accept       json
//	@Param			body 	body 		UpdateNameInput	true	 "body"
//	@Produce      json
//	@Success		200		{object}	string			"ok"
//	@Router        /user/update_name [post]
func (c *Controller) UpdateName(ctx *gin.Context) {

	userID, err := getUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := UpdateNameInput{}
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		log.Println("error ShouldBindBodyWithJSON:", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = c.storage.UpdateUserName(ctx, storeTypes.UpdateUserNameInput{
		UserID:  userID,
		NewName: req.NewName,
	})
	if err != nil {
		log.Println("storage.UpdateUserName: %w", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

type UpdateNameInput struct {
	NewName string `json:"new_name"`
}

// DeleteUser godoc
//
//	@Summary      Удаление залогиненного юзера
//	@Description  update-my-name
//	@Accept       json
//	@Produce      json
//	@Success		200		{object}	string			"ok"
//	@Router        /user/delete [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = c.usersService.CheckAndDeleteUser(ctx, usersTypes.CheckAndDeleteUserInput{
		UserID: userID,
	})
	if err != nil {
		log.Printf("usersService.CheckAndDeleteUser: %v\n", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}
