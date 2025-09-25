package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/infra/dao"
)

func NewConversationRepo(rdb rdb.Provider, idGen idgen.IDGenerator) ConversationRepo {
	return dao.NewConversationDAO(rdb.NewSession(context.Background()).DB(), idGen)
}

type ConversationRepo interface {
	Create(ctx context.Context, msg *entity.Conversation) (*entity.Conversation, error)
	GetByID(ctx context.Context, id int64) (*entity.Conversation, error)
	UpdateSection(ctx context.Context, id int64) (int64, error)
	Get(ctx context.Context, userID int64, agentID int64, scene int32) (*entity.Conversation, error)
	Update(ctx context.Context, req *entity.UpdateMeta) (*entity.Conversation, error)
	Delete(ctx context.Context, id int64) (int64, error)
	List(ctx context.Context, userID int64, agentID int64, scene int32, limit int, page int) ([]*entity.Conversation, bool, error)
}
