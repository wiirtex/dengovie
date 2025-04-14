package users

import (
	"context"
	"dengovie/internal/service/users/types"
	storeTypes "dengovie/internal/store/types"
	"fmt"
)

func (s *Service) CheckAndDeleteUser(ctx context.Context, input types.CheckAndDeleteUserInput) error {

	// TODO: транзакцию!

	debts, err := s.storage.ListUserDebts(ctx, storeTypes.ListUserDebtsInput{
		UserID: input.UserID,
	})
	if err != nil {
		return fmt.Errorf("storage.ListUserDebts: %w", err)
	}

	if len(debts) > 0 {
		return types.ErrUserHasDebts
	}

	err = s.storage.DeleteUser(ctx, storeTypes.DeleteUserInput{
		UserID: input.UserID,
	})
	if err != nil {
		return fmt.Errorf("s.storage.DeleteUser: %w", err)
	}

	return nil
}
