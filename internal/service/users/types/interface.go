package types

import "context"

type Service interface {
	CheckAndDeleteUser(ctx context.Context, input CheckAndDeleteUserInput) error
}
