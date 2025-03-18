package web

type ErrorReason string

const (
	TelegramNotFound ErrorReason = "telegram_not_found"
	InvalidOTP       ErrorReason = "invalid_otp"

	DebtorNotInGroup ErrorReason = "debtor_not_in_group"
	BuyerNotInGroup  ErrorReason = "buyer_not_in_group"
)

type APIError struct {
	ErrorReason ErrorReason
}
