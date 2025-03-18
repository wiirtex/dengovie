package types

import "errors"

var (
	ErrBuyerNotInGroup  = errors.New("buyer not in group")
	ErrDebtorNotInGroup = errors.New("debtor not in group")
)
