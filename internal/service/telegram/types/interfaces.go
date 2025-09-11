package types

import "context"

type Client interface {
	SendMessageToUserByAlias(ctx context.Context, alias string, message string) error
}
