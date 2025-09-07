package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
)

//go:generate mockgen -destination=mocks/runtime.go -package=mocks . IRuntimeRepository
type IRuntimeRepository interface {
	CreateModelRequestRecord(ctx context.Context, record *entity.ModelRequestRecord) (err error)
}
