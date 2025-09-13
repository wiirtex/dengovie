package debts

import (
	"context"
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"fmt"
)

func (s *Service) PayDebt(ctx context.Context, input debtsTypes.PayDebtInput) error {
	// TODO: транзакцию!

	debts, err := s.storage.ListUserDebts(ctx, storeTypes.ListUserDebtsInput{
		UserID: input.UserID,
	})
	if err != nil {
		return fmt.Errorf("storage.ListUserDebts: %w", err)
	}

	var debt storeTypes.UserDebt
	var found bool
	for _, d := range debts {
		if d.AnotherUser.ID == input.PayeeID {
			debt = d
			found = true
		}
	}
	if !found {
		return fmt.Errorf("debt not found")
	}

	payAmount := input.Amount
	if input.Full {
		payAmount = -debt.Amount
	} else if payAmount > (-debt.Amount) {
		return fmt.Errorf("debt amount is greater than payAmount")
	}

	if payAmount <= 0 {
		return fmt.Errorf("debt amount is non positive")
	}

	// Обновляем записи о долгах
	err = s.storage.ShareDebt(ctx, storeTypes.ShareDebtInput{
		UserID: input.PayeeID,
		ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
			{
				UserID: input.UserID,
				Amount: -payAmount,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("storage.ShareDebt: %w", err)
	}

	return nil

}
