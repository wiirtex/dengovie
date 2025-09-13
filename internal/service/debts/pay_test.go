package debts

import (
	"context"
	"dengovie/internal/domain"
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_PayDebt(t *testing.T) {
	type args struct {
		input debtsTypes.PayDebtInput
	}
	tests := []struct {
		name         string
		args         args
		prepareMocks func(e *env)
		wantErr      bool
		errContains  string
	}{
		{
			name: "OK - successful partial payment",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  500,
				},
			},
			prepareMocks: func(e *env) {
				// Mock list debts
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000, // Долг 1000
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
					{
						Amount: -500, // Другой долг
						AnotherUser: storeTypes.User{
							ID: 789,
						},
					},
				}, nil)

				// Mock share debt (обновление долга)
				e.mockStorage.EXPECT().ShareDebt(mock.Anything, storeTypes.ShareDebtInput{
					UserID: domain.UserID(456),
					ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
						{
							UserID: domain.UserID(123),
							Amount: -500, // Уменьшаем долг на 500
						},
					},
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "OK - successful full payment",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    true,
					Amount:  0, // Amount игнорируется при Full=true
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1500, // Долг 1500
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
				}, nil)

				e.mockStorage.EXPECT().ShareDebt(mock.Anything, storeTypes.ShareDebtInput{
					UserID: domain.UserID(456),
					ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
						{
							UserID: domain.UserID(123),
							Amount: -1500, // Полное погашение долга
						},
					},
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "FAIL - debt not found",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  500,
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000,
						AnotherUser: storeTypes.User{
							ID: 789, // Другой пользователь
						},
					},
				}, nil)
			},
			wantErr:     true,
			errContains: "debt not found",
		},
		{
			name: "FAIL - storage error listing debts",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  500,
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return(nil, assert.AnError)
			},
			wantErr:     true,
			errContains: "storage.ListUserDebts",
		},
		{
			name: "FAIL - payment amount exceeds debt",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  2000, // Больше чем долг
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000, // Долг только 1000
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
				}, nil)
			},
			wantErr:     true,
			errContains: "debt amount is greater than payAmount",
		},
		{
			name: "FAIL - negative payment amount",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  -100, // Отрицательная сумма
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000,
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
				}, nil)
			},
			wantErr:     true,
			errContains: "debt amount is non positive",
		},
		{
			name: "FAIL - zero payment amount without full flag",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  0, // Нулевая сумма без full
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000,
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
				}, nil)
			},
			wantErr:     true,
			errContains: "debt amount is non positive", // Так как 0 < 0 = false, но debt отрицательный
		},
		{
			name: "FAIL - storage error when updating debt",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  500,
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000,
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
				}, nil)

				e.mockStorage.EXPECT().ShareDebt(mock.Anything, storeTypes.ShareDebtInput{
					UserID: domain.UserID(456),
					ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
						{
							UserID: domain.UserID(123),
							Amount: -500,
						},
					},
				}).Return(assert.AnError)
			},
			wantErr:     true,
			errContains: "storage.ShareDebt",
		},
		{
			name: "OK - pay exact debt amount",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  1000, // Точная сумма долга
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -1000,
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
				}, nil)

				e.mockStorage.EXPECT().ShareDebt(mock.Anything, storeTypes.ShareDebtInput{
					UserID: domain.UserID(456),
					ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
						{
							UserID: domain.UserID(123),
							Amount: -1000, // Полное погашение
						},
					},
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "OK - multiple debts, find correct one",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  300,
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{
					{
						Amount: -500,
						AnotherUser: storeTypes.User{
							ID: 789, // Первый долг (не тот)
						},
					},
					{
						Amount: -1000, // Второй долг (правильный)
						AnotherUser: storeTypes.User{
							ID: 456,
						},
					},
					{
						Amount: -200, // Третий долг (не тот)
						AnotherUser: storeTypes.User{
							ID: 999,
						},
					},
				}, nil)

				e.mockStorage.EXPECT().ShareDebt(mock.Anything, storeTypes.ShareDebtInput{
					UserID: domain.UserID(456),
					ChangeDebtAmount: []storeTypes.ChangeUserDebtAmountInput{
						{
							UserID: domain.UserID(123),
							Amount: -300,
						},
					},
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "FAIL - no debts at all",
			args: args{
				input: debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  500,
				},
			},
			prepareMocks: func(e *env) {
				e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.UserDebt{}, nil)
			},
			wantErr:     true,
			errContains: "debt not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			s := newService(e)

			tt.prepareMocks(e)

			err := s.PayDebt(context.Background(), tt.args.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}

			e.mockStorage.AssertExpectations(t)
		})
	}
}

// Дополнительные тесты для edge cases
func TestService_PayDebt_EdgeCases(t *testing.T) {
	t.Run("zero debt amount", func(t *testing.T) {
		e := newEnv(t)
		s := newService(e)

		e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
			UserID: domain.UserID(123),
		}).Return([]storeTypes.UserDebt{
			{
				Amount: 0, // Нулевой долг
				AnotherUser: storeTypes.User{
					ID: 456,
				},
			},
		}, nil)

		input := debtsTypes.PayDebtInput{
			UserID:  domain.UserID(123),
			PayeeID: domain.UserID(456),
			Full:    false,
			Amount:  100,
		}

		err := s.PayDebt(context.Background(), input)
		require.Error(t, err)
		require.Contains(t, err.Error(), "debt amount is greater than payAmount")
	})

	t.Run("positive debt amount (credit)", func(t *testing.T) {
		e := newEnv(t)
		s := newService(e)

		e.mockStorage.EXPECT().ListUserDebts(mock.Anything, storeTypes.ListUserDebtsInput{
			UserID: domain.UserID(123),
		}).Return([]storeTypes.UserDebt{
			{
				Amount: 500, // Положительный долг (кредит)
				AnotherUser: storeTypes.User{
					ID: 456,
				},
			},
		}, nil)

		input := debtsTypes.PayDebtInput{
			UserID:  domain.UserID(123),
			PayeeID: domain.UserID(456),
			Full:    false,
			Amount:  100,
		}

		err := s.PayDebt(context.Background(), input)
		require.Error(t, err)
		require.Contains(t, err.Error(), "debt amount is greater than payAmount")
	})
}
