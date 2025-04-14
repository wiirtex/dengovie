package types

import "dengovie/internal/domain"

type ShareDebtInput struct {
	BuyerID   domain.UserID   `json:"buyer_id"`
	GroupID   domain.GroupID  `json:"group_id"`
	DebtorIDs []domain.UserID `json:"debtor_ids"`
	Amount    int64           `json:"amount"`
}

type PayDebtInput struct {
	UserID  domain.UserID `json:"payer_user_id"`
	PayeeID domain.UserID `json:"payee_user_id"`
	Full    bool          `json:"full"`
	Amount  int64         `json:"amount"`
}
