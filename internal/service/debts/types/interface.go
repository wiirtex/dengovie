package types

import "context"

type DebtsService interface {
	ShareDebt(ctx context.Context, input ShareDebtInput) error
	PayDebt(ctx context.Context, input PayDebtInput) error
}
