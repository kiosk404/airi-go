package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	message "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/infra/repo/dao"
)

func NewMessageRepo(rdb rdb.Provider, idGen idgen.IDGenerator) MessageRepo {
	return dao.NewMessageDAO(rdb.NewSession(context.Background()).DB(), idGen)
}

type MessageRepo interface {
	PreCreate(ctx context.Context, msg *entity.Message) (*entity.Message, error)
	Create(ctx context.Context, msg *entity.Message) (*entity.Message, error)
	BatchCreate(ctx context.Context, msg []*entity.Message) ([]*entity.Message, error)
	List(ctx context.Context, listMeta *entity.ListMeta) ([]*entity.Message, bool, error)
	GetByRunIDs(ctx context.Context, runIDs []int64, orderBy string) ([]*entity.Message, error)
	Edit(ctx context.Context, msgID int64, message *message.Message) (int64, error)
	GetByID(ctx context.Context, msgID int64) (*entity.Message, error)
	Delete(ctx context.Context, delMeta *entity.DeleteMeta) error
}
