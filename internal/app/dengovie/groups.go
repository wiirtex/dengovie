package dengovie

import (
	"context"
	"dengovie/internal/domain"
	storeTypes "dengovie/internal/store/types"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListUserGroups godoc
//
//	@Summary      Вернуть список групп пользователя, в которых он состоит
//	@Description  list-user-groups
//	@Accept       json
//	@Produce      json
//	@Success		200		{object}	string			"ok"
//	@Failure		400		{object}	web.APIError	"невалидный запрос"
//	@Failure		404		{object}	web.APIError	"клиент не найден"
//	@Router        /groups [get]
func (c *Controller) ListUserGroups(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	groups, err := c.storage.ListUserGroups(context.TODO(), storeTypes.ListUserGroupsInput{
		UserID: domain.UserID(userID),
	})
	if err != nil {
		log.Println("storage.ListUserGroups:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, groups)
}

// ListUsersInGroup godoc
//
//	@Summary      Список юзеров в группе
//	@Description  list-users-in-group
//	@Accept       json
//	@Produce      json
//	@Param			groupID 	path 	int64	true	 "groupID"
//	@Success		200		{object}	string			"ok"
//	@Failure		400		{object}	web.APIError	"невалидный запрос"
//	@Failure		404		{object}	web.APIError	"группа не найдена"
//	@Router        /groups/{groupID}/users [get]
func (c *Controller) ListUsersInGroup(ctx *gin.Context) {

	_, err := getUserID(ctx) // тут будет защита от IDOR
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(ctx.Param("groupID"))
	if err != nil {
		log.Println("strconv.Atoi(groupIDStr):", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	groups, err := c.storage.ListUsersInGroup(context.TODO(), storeTypes.ListUsersInGroupInput{
		GroupID: domain.GroupID(groupID),
	})
	if err != nil {
		log.Println("storage.ListUsersInGroup:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, groups)
}
