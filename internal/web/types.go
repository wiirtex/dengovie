package web

type ErrorReason string

const (
	TelegramNotFound ErrorReason = "telegram_not_found"
	InvalidOTP       ErrorReason = "invalid_otp"
)

type APIError struct {
	ErrorReason ErrorReason
}
