package debts

import (
	"context"
	"dengovie/internal/domain"
	"dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"fmt"
	"math/rand/v2"
)

func (s *Service) ShareDebt(ctx context.Context, input types.ShareDebtInput) error {

	// TODO: транзакцию!

	// проверяем, что деляший в группе
	isIn, err := s.storage.IsUserInGroup(ctx, input.BuyerID, input.GroupID)
	if err != nil {
		return fmt.Errorf("storage.IsUserInGroup: %w", err)
	}
	if !isIn {
		return types.ErrBuyerNotInGroup
	}

	// проверяем, что все должники в группе
	areIn, err := s.storage.AreUsersInGroup(ctx, input.DebtorIDs, input.GroupID)
	if err != nil {
		return fmt.Errorf("storage.AreUsersInGroup: %w", err)
	}
	if !areIn {
		return types.ErrDebtorNotInGroup
	}

	accounts := make([]domain.UserID, 0, len(input.DebtorIDs))
	for _, id := range input.DebtorIDs {
		if id == input.BuyerID {
			continue
		}

		accounts = append(accounts, id)
	}

	// Создаём пустые записи о долгах между людьми в группе
	err = s.storage.CreateEmptyDebts(ctx, storeTypes.CreateEmptyDebtsInput{
		UserID:         input.BuyerID,
		AnotherUserIDs: accounts,
	})
	if err != nil {
		return fmt.Errorf("storage.CreateEmptyDebts: %w", err)
	}

	// Сумма делится строго без остатка. Первые remainder дебторов получат +1 копейку в случае, если деление не целое.
	// Так что шаффлим их, чтобы каждый раз разные дебторы получали копейку
	rand.Shuffle(len(input.DebtorIDs), func(i, j int) {
		input.DebtorIDs[i], input.DebtorIDs[j] = input.DebtorIDs[j], input.DebtorIDs[i]
	})

	sum := input.Amount / int64(len(input.DebtorIDs))
	remainder := input.Amount % int64(len(input.DebtorIDs))
	amounts := make([]storeTypes.ChangeUserDebtAmountInput, 0, len(input.DebtorIDs))
	for _, id := range input.DebtorIDs {
		diff := int64(0)
		if remainder > 0 {
			diff = 1
			remainder--
		}

		if id == input.BuyerID {
			continue
		}

		amounts = append(amounts, storeTypes.ChangeUserDebtAmountInput{
			UserID: id,
			Amount: sum + diff,
		})
	}

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
