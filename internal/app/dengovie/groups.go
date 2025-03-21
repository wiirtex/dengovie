package dengovie

import (
	"context"
	"dengovie/internal/domain"
	storeTypes "dengovie/internal/store/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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

	userID, err := strconv.Atoi(ctx.GetString(domain.UserIDKey))
	if err != nil {
		log.Println("strconv.Atoi:", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
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

	_ = ctx.GetInt64(domain.UserIDKey) // тут будет защита от IDOR
	groupID, err := strconv.Atoi(ctx.Param("groupID"))
	if err != nil {
		log.Println("strconv.Atoi(groupIDStr):", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	groups, err := c.storage.ListUsersInGroup(context.TODO(), storeTypes.ListUsersInGroupInput{
		GroupID: domain.GroupID(groupID),
	})
	if err != nil {
		log.Println("storage.ListUsersInGroup:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, groups)
}
