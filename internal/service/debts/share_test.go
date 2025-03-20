package debts

import (
	"context"
	"dengovie/internal/domain"
	"dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_ShareDebt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        types.ShareDebtInput
		prepareMocks func(t *testing.T, e *env)
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			input: types.ShareDebtInput{
				BuyerID: 1,
				GroupID: 2,
				DebtorIDs: []domain.UserID{
					3, 4, 5,
				},
				Amount: 6,
			},
			wantErr: assert.NoError,
			prepareMocks: func(t *testing.T, e *env) {

				e.mockStorage.EXPECT().
					IsUserInGroup(mock.Anything, domain.UserID(1), domain.GroupID(2)).
					Return(true, nil)

				e.mockStorage.EXPECT().
					AreUsersInGroup(mock.Anything, []domain.UserID{3, 4, 5}, domain.GroupID(2)).
					Return(true, nil)

				e.mockStorage.EXPECT().
					CreateEmptyDebts(mock.Anything, storeTypes.CreateEmptyDebtsInput{
						UserID:         1,
						AnotherUserIDs: []domain.UserID{3, 4, 5},
					}).
					Return(nil)

				e.mockStorage.EXPECT().
					ShareDebt(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input storeTypes.ShareDebtInput) error {

						diff := cmp.Diff(storeTypes.ShareDebtInput{
							UserID: 1,
							ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
								{UserID: 3},
								{UserID: 5},
								{UserID: 4},
							},
						}, input,
							cmpopts.IgnoreFields(storeTypes.ChangeUserDebtAmountInput{}, "Amount"),
							cmpopts.SortSlices(func(a, b storeTypes.ChangeUserDebtAmountInput) bool {
								return a.UserID < b.UserID
							}))
						assert.Empty(t, diff, "все передаваемые ID должны быть равны")

						sum := int64(0)
						for _, u := range input.ChangeDebtAmount {
							sum += u.Amount
						}
						assert.Equal(t, int64(6), sum, "сумма всех частичных долгов должна быть равна переданной")

						amounts := input.ChangeDebtAmount
						for _, i := range amounts {
							for _, j := range amounts {
								value := max(i.Amount, j.Amount)-min(j.Amount, i.Amount) <= 1
								assert.True(t, value, "попарно ВСЕ суммы не должны отличаться больше, чем на 1")
							}
						}

						return nil
					})
			},
		},
		{
			name: "Fail, CreateEmptyDebts",
			input: types.ShareDebtInput{
				BuyerID: 1,
				GroupID: 2,
				DebtorIDs: []domain.UserID{
					3, 4, 5,
				},
				Amount: 6,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, assert.AnError)
			},
			prepareMocks: func(t *testing.T, e *env) {

				e.mockStorage.EXPECT().
					IsUserInGroup(mock.Anything, domain.UserID(1), domain.GroupID(2)).
					Return(true, nil)

				e.mockStorage.EXPECT().
					AreUsersInGroup(mock.Anything, []domain.UserID{3, 4, 5}, domain.GroupID(2)).
					Return(true, nil)

				e.mockStorage.EXPECT().
					CreateEmptyDebts(mock.Anything, storeTypes.CreateEmptyDebtsInput{
						UserID:         1,
						AnotherUserIDs: []domain.UserID{3, 4, 5},
					}).
					Return(assert.AnError)
			},
		},
		{
			name: "Fail, AreUsersInGroup",
			input: types.ShareDebtInput{
				BuyerID: 1,
				GroupID: 2,
				DebtorIDs: []domain.UserID{
					3, 4, 5,
				},
				Amount: 6,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, assert.AnError)
			},
			prepareMocks: func(t *testing.T, e *env) {

				e.mockStorage.EXPECT().
					IsUserInGroup(mock.Anything, domain.UserID(1), domain.GroupID(2)).
					Return(true, nil)

				e.mockStorage.EXPECT().
					AreUsersInGroup(mock.Anything, []domain.UserID{3, 4, 5}, domain.GroupID(2)).
					Return(true, assert.AnError)
			},
		},
		{
			name: "Fail, debtors not in group",
			input: types.ShareDebtInput{
				BuyerID: 1,
				GroupID: 2,
				DebtorIDs: []domain.UserID{
					3, 4, 5,
				},
				Amount: 6,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, types.ErrDebtorNotInGroup)
			},
			prepareMocks: func(t *testing.T, e *env) {

				e.mockStorage.EXPECT().
					IsUserInGroup(mock.Anything, domain.UserID(1), domain.GroupID(2)).
					Return(true, nil)

				e.mockStorage.EXPECT().
					AreUsersInGroup(mock.Anything, []domain.UserID{3, 4, 5}, domain.GroupID(2)).
					Return(false, nil)
			},
		},
		{
			name: "Fail, IsUserInGroup",
			input: types.ShareDebtInput{
				BuyerID: 1,
				GroupID: 2,
				DebtorIDs: []domain.UserID{
					3, 4, 5,
				},
				Amount: 6,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, assert.AnError)
			},
			prepareMocks: func(t *testing.T, e *env) {

				e.mockStorage.EXPECT().
					IsUserInGroup(mock.Anything, domain.UserID(1), domain.GroupID(2)).
					Return(true, assert.AnError)
			},
		},
		{
			name: "Fail, buyer not in group",
			input: types.ShareDebtInput{
				BuyerID: 1,
				GroupID: 2,
				DebtorIDs: []domain.UserID{
					3, 4, 5,
				},
				Amount: 6,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, types.ErrBuyerNotInGroup)
			},
			prepareMocks: func(t *testing.T, e *env) {

				e.mockStorage.EXPECT().
					IsUserInGroup(mock.Anything, domain.UserID(1), domain.GroupID(2)).
					Return(false, nil)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := newEnv(t)
			tt.prepareMocks(t, e)

			s := newService(e)
			err := s.ShareDebt(context.Background(), tt.input)

			assert.True(t, tt.wantErr(t, err))
		})
	}
}
