package rpc

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/component/rpc"
)

type AuthRPCAdapter struct {
}

func NewAuthRPCProvider() rpc.IAuthProvider {
	return &AuthRPCAdapter{}
}

func (a *AuthRPCAdapter) CheckPermission(ctx context.Context, action string) error {
	// TODO: implement
	return nil
}
