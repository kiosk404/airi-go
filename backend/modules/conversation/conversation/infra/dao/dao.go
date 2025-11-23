package dao

import (
	"context"
	"errors"
	"time"

	"github.com/kiosk404/airi-go/backend/api/model/conversation/common"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/infra/repo/gorm_gen/query"
	model2 "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/conversation/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"gorm.io/gorm"
)

type ConversationDAO struct {
	idGen idgen.IDGenerator
	db    *gorm.DB
	query *query.Query
}

func NewConversationDAO(db *gorm.DB, generator idgen.IDGenerator) *ConversationDAO {
	return &ConversationDAO{
		idGen: generator,
		db:    db,
		query: query.Use(db),
	}
}

func (dao *ConversationDAO) Create(ctx context.Context, msg *entity.Conversation) (*entity.Conversation, error) {
	poData := dao.conversationDO2PO(ctx, msg)

	ids, err := dao.idGen.GenMultiIDs(ctx, 2)
	if err != nil {
		return nil, err
	}
	poData.ID = ids[0]
	poData.SectionID = ids[1]

	err = dao.query.Conversation.WithContext(ctx).Create(poData)
	if err != nil {
		return nil, err
	}
	return dao.conversationPO2DO(ctx, poData), nil
}

func (dao *ConversationDAO) GetByID(ctx context.Context, id int64) (*entity.Conversation, error) {
	poData, err := dao.query.Conversation.WithContext(ctx).Debug().Where(dao.query.Conversation.ID.Eq(id)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return dao.conversationPO2DO(ctx, poData), nil
}

func (dao *ConversationDAO) UpdateSection(ctx context.Context, id int64) (int64, error) {
	updateColumn := make(map[string]interface{})
	table := dao.query.Conversation
	newSectionID, err := dao.idGen.GenID(ctx)
	if err != nil {
		return 0, err
	}
	updateColumn[table.SectionID.ColumnName().String()] = newSectionID
	updateColumn[table.UpdatedAt.ColumnName().String()] = time.Now().UnixMilli()

	_, err = dao.query.Conversation.WithContext(ctx).Where(dao.query.Conversation.ID.Eq(id)).UpdateColumns(updateColumn)
	if err != nil {
		return 0, err
	}
	return newSectionID, nil
}

func (dao *ConversationDAO) Delete(ctx context.Context, id int64) (int64, error) {
	table := dao.query.Conversation

	updateColumn := make(map[string]interface{})
	updateColumn[table.UpdatedAt.ColumnName().String()] = time.Now().UnixMilli()
	updateColumn[table.Status.ColumnName().String()] = model2.ConversationStatusDeleted

	updateRes, err := dao.query.Conversation.WithContext(ctx).Where(dao.query.Conversation.ID.Eq(id)).UpdateColumns(updateColumn)
	if err != nil {
		return 0, err
	}
	return updateRes.RowsAffected, err
}

func (dao *ConversationDAO) Update(ctx context.Context, req *entity.UpdateMeta) (*entity.Conversation, error) {
	updateColumn := make(map[string]interface{})
	updateColumn[dao.query.Conversation.UpdatedAt.ColumnName().String()] = time.Now().UnixMilli()
	if len(req.Name) > 0 {
		updateColumn[dao.query.Conversation.Name.ColumnName().String()] = req.Name
	}

	_, err := dao.query.Conversation.WithContext(ctx).Where(dao.query.Conversation.ID.Eq(req.ID)).UpdateColumns(updateColumn)
	if err != nil {
		return nil, err
	}
	return dao.GetByID(ctx, req.ID)
}

func (dao *ConversationDAO) Get(ctx context.Context, userID int64, agentID int64, scene int32) (*entity.Conversation, error) {
	po, err := dao.query.Conversation.WithContext(ctx).Debug().
		Where(dao.query.Conversation.CreatorID.Eq(userID)).
		Where(dao.query.Conversation.AgentID.Eq(agentID)).
		Where(dao.query.Conversation.Scene.Eq(scene)).
		Where(dao.query.Conversation.Status.Eq(int32(model2.ConversationStatusNormal))).
		Order(dao.query.Conversation.CreatedAt.Desc()).
		First()

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dao.conversationPO2DO(ctx, po), nil
}

func (dao *ConversationDAO) List(ctx context.Context, userID int64, agentID int64, scene int32, limit int, page int) ([]*entity.Conversation, bool, error) {
	var hasMore bool

	do := dao.query.Conversation.WithContext(ctx).Debug()
	do = do.Where(dao.query.Conversation.CreatorID.Eq(userID)).
		Where(dao.query.Conversation.AgentID.Eq(agentID)).
		Where(dao.query.Conversation.Scene.Eq(scene)).
		Where(dao.query.Conversation.Status.Eq(int32(model2.ConversationStatusNormal)))

	do = do.Offset((page - 1) * limit)

	if limit > 0 {
		do = do.Limit(int(limit) + 1)
	}
	do = do.Order(dao.query.Conversation.CreatedAt.Desc())

	poList, err := do.Find()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, hasMore, nil
	}
	if err != nil {
		return nil, hasMore, err
	}

	if len(poList) == 0 {
		return nil, hasMore, nil
	}
	if len(poList) > limit {
		hasMore = true
		return dao.conversationBatchPO2DO(ctx, poList[:(len(poList)-1)]), hasMore, nil

	}
	return dao.conversationBatchPO2DO(ctx, poList), hasMore, nil
}

func (dao *ConversationDAO) conversationDO2PO(ctx context.Context, conversation *entity.Conversation) *model.Conversation {
	return &model.Conversation{
		ID:        conversation.ID,
		SectionID: conversation.SectionID,
		AgentID:   conversation.AgentID,
		CreatorID: ptr.Of(conversation.CreatorID),
		Scene:     int32(conversation.Scene),
		Status:    int32(conversation.Status),
		Ext:       ptr.Of(conversation.Ext),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
		Name:      conversation.Name,
	}
}

func (dao *ConversationDAO) conversationPO2DO(ctx context.Context, c *model.Conversation) *entity.Conversation {
	return &entity.Conversation{
		ID:        c.ID,
		SectionID: c.SectionID,
		AgentID:   c.AgentID,
		CreatorID: ptr.From(c.CreatorID),
		Scene:     common.Scene(c.Scene),
		Status:    model2.ConversationStatus(c.Status),
		Ext:       ptr.From(c.Ext),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}

func (dao *ConversationDAO) conversationBatchPO2DO(ctx context.Context, conversations []*model.Conversation) []*entity.Conversation {
	return slices.Transform(conversations, func(c *model.Conversation) *entity.Conversation {
		return &entity.Conversation{
			ID:        c.ID,
			SectionID: c.SectionID,
			AgentID:   c.AgentID,
			CreatorID: ptr.From(c.CreatorID),
			Scene:     common.Scene(c.Scene),
			Status:    model2.ConversationStatus(c.Status),
			Ext:       ptr.From(c.Ext),
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Name:      c.Name,
		}
	})
}
