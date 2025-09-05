package rpc

import (
	"context"
)

//go:generate mockgen -destination=mocks/auth_provider.go -package=mocks . IAuthProvider
type IAuthProvider interface {
	CheckPermission(ctx context.Context, action string) error
}
