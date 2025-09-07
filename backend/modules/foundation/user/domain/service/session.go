package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
)

//go:generate mockgen -destination=mocks/session_service.go -package=mock_session . ISessionService
type ISessionService interface {
	ValidateSession(ctx context.Context, sessionID string) (*entity.Session, error)
	GenerateSessionKey(ctx context.Context, session *entity.Session) (string, error)
}
