package types

import "context"

type DebtsService interface {
	ShareDebt(ctx context.Context, input ShareDebtInput) error
}
