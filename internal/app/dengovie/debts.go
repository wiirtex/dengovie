package dengovie

import (
	"context"
	"dengovie/internal/domain"
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"dengovie/internal/web"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// ListDebts godoc
//
//	@Summary      Список всех долгов юзера, в которых он состоит
//	@Description  list-user-groups
//	@Accept       json
//	@Produce      json
//	@Success		200		{object}	ListDebtsResponseBody			"ok"
//	@Failure		400		{object}	web.APIError	"невалидный запрос"
//	@Failure		404		{object}	web.APIError	"клиент не найден"
//	@Router        /debts [get]
func (c *Controller) ListDebts(ctx *gin.Context) {

	userID := ctx.GetString(domain.UserIDKey)
	id, err := strconv.Atoi(userID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	debts, err := c.storage.ListUserDebts(context.TODO(), storeTypes.ListUserDebtsInput{
		UserID: domain.UserID(id),
	})
	if err != nil {
		log.Println("storage.ListUserDebts:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var resp ListDebtsResponseBody
	for _, debt := range debts {
		resp.Debts = append(resp.Debts, UserDebt{
			AnotherUserID:   debt.AnotherUser.ID,
			AnotherUserName: debt.AnotherUser.Name,
			Amount:          debt.Amount,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

type ListDebtsResponseBody struct {
	Debts []UserDebt
}

type UserDebt struct {
	AnotherUserID   domain.UserID `json:"another_user_id"`
	AnotherUserName string        `json:"another_user_name"`
	Amount          int64         `json:"amount"`
}

// ShareDebt godoc
//
//	@Summary      Разделить долг между пользователями
//	@Description  list-user-groups
//	@Accept       json
//	@Produce      json
//	@Param			body 	body 	ShareDebtRequest	true	 "body"
//	@Failure		400		{object}	web.APIError	"невалидный запрос"
//	@Router        /debts/share [post]
func (c *Controller) ShareDebt(ctx *gin.Context) {

	userID, err := strconv.Atoi(ctx.GetString(domain.UserIDKey))
	if err != nil {
		log.Println("strconv.Atoi(userID):", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := ShareDebtRequest{}
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		log.Println("error ShouldBindBodyWithJSON:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.debtsService.ShareDebt(context.TODO(), debtsTypes.ShareDebtInput{
		BuyerID:   domain.UserID(userID),
		GroupID:   req.GroupID,
		DebtorIDs: req.UserIDs,
		Amount:    req.Amount,
	})
	if err != nil {
		if errors.Is(err, debtsTypes.ErrDebtorNotInGroup) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, web.APIError{
				ErrorReason: web.DebtorNotInGroup,
			})
			return
		}
		if errors.Is(err, debtsTypes.ErrBuyerNotInGroup) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, web.APIError{
				ErrorReason: web.BuyerNotInGroup,
			})
			return
		}

		log.Println("debtsService.ShareDebt:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

type ShareDebtRequest struct {
	GroupID domain.GroupID  `json:"group_id"`
	UserIDs []domain.UserID `json:"user_ids"`
	Amount  int64           `json:"amount"`
}
