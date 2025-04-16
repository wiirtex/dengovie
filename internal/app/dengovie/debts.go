package dengovie

import (
	"context"
	"dengovie/internal/domain"
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"dengovie/internal/web"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	debts, err := c.storage.ListUserDebts(context.TODO(), storeTypes.ListUserDebtsInput{
		UserID: userID,
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
	userID, err := getUserID(ctx)
	if err != nil {
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

// PayDebt godoc
//
//	@Summary      	Выплатить долг пользователю
//	@Description  	pay-debt-to-user
//	@Accept       	json
//	@Produce    	json
//	@Param			body 	body 		PayDebtRequest	true	 "body"
//	@Failure		400		{object}	web.APIError		"невалидный запрос"
//	@Router        	/debts/pay [post]
func (c *Controller) PayDebt(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := PayDebtRequest{}
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		log.Println("error ShouldBindBodyWithJSON:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.debtsService.PayDebt(context.TODO(), debtsTypes.PayDebtInput{
		UserID:  domain.UserID(userID),
		PayeeID: req.AnotherUserID,
		Full:    req.Full,
		Amount:  req.Amount,
	})
	if err != nil {
		// TODO: добавить сюда типизированные ошибки
		log.Println("debtsService.ShareDebt:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

type PayDebtRequest struct {
	// Долг перед кем выплачивается
	AnotherUserID domain.UserID `json:"another_user_id"`
	// Full если true, то сумма не учитывается
	Full   bool  `json:"full"`
	Amount int64 `json:"amount"`
}
