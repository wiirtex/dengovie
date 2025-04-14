package debts

import (
	"context"
	"dengovie/internal/domain"
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"fmt"
)

func (s *Service) PayDebt(ctx context.Context, input debtsTypes.PayDebtInput) error {
	// TODO: транзакцию!

	// Обновляем записи о долгах
	err = s.storage.ShareDebt(ctx, storeTypes.ShareDebtInput{
		UserID:           input.BuyerID,
		ChangeDebtAmount: amounts,
	})
	if err != nil {
		return fmt.Errorf("storage.ShareDebt: %w", err)
	}

	return nil

}
