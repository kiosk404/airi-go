package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/crossdomain/message"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/infra/repo/dao"
	"gorm.io/gorm"
)

func NewMessageRepo(db *gorm.DB, idGen idgen.IDGenerator) MessageRepo {
	return dao.NewMessageDAO(db, idGen)
}

type MessageRepo interface {
	PreCreate(ctx context.Context, msg *entity.Message) (*entity.Message, error)
	Create(ctx context.Context, msg *entity.Message) (*entity.Message, error)
	List(ctx context.Context, listMeta *entity.ListMeta) ([]*entity.Message, bool, error)
	GetByRunIDs(ctx context.Context, runIDs []int64, orderBy string) ([]*entity.Message, error)
	Edit(ctx context.Context, msgID int64, message *message.Message) (int64, error)
	GetByID(ctx context.Context, msgID int64) (*entity.Message, error)
	Delete(ctx context.Context, delMeta *entity.DeleteMeta) error
}
