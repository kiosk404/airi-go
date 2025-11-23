package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/conversation/common"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/conversation"
	ctxutil2 "github.com/kiosk404/airi-go/backend/application/ctxutil"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/entity"
	conversationService "github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	message "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/service"
	uploadService "github.com/kiosk404/airi-go/backend/modules/data/upload/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

type ConversationApplicationService struct {
	appContext *ServiceComponents

	AgentRunDomainSVC     agentrun.Run
	ConversationDomainSVC conversationService.Conversation
	MessageDomainSVC      message.Message
}

var ConversationSVC = new(ConversationApplicationService)

func (c *ConversationApplicationService) ClearHistory(ctx context.Context, req *conversation.ClearConversationHistoryRequest) (*conversation.ClearConversationHistoryResponse, error) {
	resp := new(conversation.ClearConversationHistoryResponse)

	conversationID := req.ConversationID

	// get conversation
	currentRes, err := c.ConversationDomainSVC.GetByID(ctx, conversationID)
	if err != nil {
		return resp, err
	}
	if currentRes == nil {
		return resp, errorx.New(errno.ErrConversationNotFound)
	}
	// check user
	userID := ctxutil2.GetUIDFromCtx(ctx)
	if userID == nil || *userID != currentRes.CreatorID {
		return resp, errorx.New(errno.ErrConversationNotFound, errorx.KV("msg", "user not match"))
	}
	// delete conversation
	err = c.ConversationDomainSVC.Delete(ctx, conversationID)
	if err != nil {
		return resp, err
	}
	// create new conversation
	convRes, err := c.ConversationDomainSVC.Create(ctx, &entity.CreateMeta{
		AgentID: currentRes.AgentID,
		UserID:  currentRes.CreatorID,
		Scene:   currentRes.Scene,
	})
	if err != nil {
		return resp, err
	}
	resp.NewSectionID = convRes.SectionID
	return resp, nil
}

func (c *ConversationApplicationService) CreateSection(ctx context.Context, conversationID int64) (int64, error) {
	currentRes, err := c.ConversationDomainSVC.GetByID(ctx, conversationID)
	if err != nil {
		return 0, err
	}

	if currentRes == nil {
		return 0, errorx.New(errno.ErrConversationNotFound, errorx.KV("msg", "conversation not found"))
	}
	var userID int64

	userID = ctxutil2.MustGetUIDFromCtx(ctx)

	if userID != currentRes.CreatorID {
		return 0, errorx.New(errno.ErrConversationNotFound, errorx.KV("msg", "user not match"))
	}

	convRes, err := c.ConversationDomainSVC.NewConversationCtx(ctx, &entity.NewConversationCtxRequest{
		ID: conversationID,
	})
	if err != nil {
		return 0, err
	}
	return convRes.SectionID, nil
}

func (c *ConversationApplicationService) CreateConversation(ctx context.Context, req *conversation.CreateConversationRequest) (*conversation.CreateConversationResponse, error) {
	resp := new(conversation.CreateConversationResponse)
	apiKeyInfo := ctxutil2.GetApiAuthFromCtx(ctx)
	userID := apiKeyInfo.UserID
	agentID := req.GetBotId()

	conversationData, err := c.ConversationDomainSVC.Create(ctx, &entity.CreateMeta{
		AgentID: agentID,
		UserID:  userID,
		Scene:   common.Scene_SceneOpenApi,
		Ext:     parseMetaData(req.MetaData),
	})
	if err != nil {
		return nil, err
	}
	resp.ConversationData = &conversation.ConversationData{
		Id:            conversationData.ID,
		LastSectionID: &conversationData.SectionID,
		CreatedAt:     conversationData.CreatedAt / 1000,
		MetaData:      parseExt(conversationData.Ext),
	}
	return resp, nil
}

func (c *ConversationApplicationService) ListConversation(ctx context.Context, req *conversation.ListConversationsApiRequest) (*conversation.ListConversationsApiResponse, error) {

	resp := new(conversation.ListConversationsApiResponse)

	apiKeyInfo := ctxutil2.GetApiAuthFromCtx(ctx)
	userID := apiKeyInfo.UserID

	if userID == 0 {
		return resp, errorx.New(errno.ErrConversationNotFound)
	}

	conversationDOList, hasMore, err := c.ConversationDomainSVC.List(ctx, &entity.ListMeta{
		UserID:  userID,
		AgentID: req.GetBotID(),
		Scene:   common.Scene_SceneOpenApi,
		Page:    int(req.GetPageNum()),
		Limit:   int(req.GetPageSize()),
	})
	if err != nil {
		return resp, err
	}
	conversationData := slices.Transform(conversationDOList, func(conv *entity.Conversation) *conversation.ConversationData {
		return &conversation.ConversationData{
			Id:            conv.ID,
			LastSectionID: &conv.SectionID,
			CreatedAt:     conv.CreatedAt / 1000,
			Name:          ptr.Of(conv.Name),
			MetaData:      parseExt(conv.Ext),
		}
	})

	resp.Data = &conversation.ListConversationData{
		Conversations: conversationData,
		HasMore:       hasMore,
	}
	return resp, nil
}

func (c *ConversationApplicationService) DeleteConversation(ctx context.Context, req *conversation.DeleteConversationApiRequest) (*conversation.DeleteConversationApiResponse, error) {
	resp := new(conversation.DeleteConversationApiResponse)
	convID := req.GetConversationID()

	apiKeyInfo := ctxutil2.GetApiAuthFromCtx(ctx)
	userID := apiKeyInfo.UserID

	if userID == 0 {
		return resp, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "permission check failed"))
	}

	conversationDO, err := c.ConversationDomainSVC.GetByID(ctx, convID)
	if err != nil {
		return resp, err
	}
	if conversationDO == nil {
		return resp, errorx.New(errno.ErrConversationNotFound)
	}
	if conversationDO.CreatorID != userID {
		return resp, errorx.New(errno.ErrConversationNotFound, errorx.KV("msg", "user not match"))
	}
	err = c.ConversationDomainSVC.Delete(ctx, convID)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *ConversationApplicationService) UpdateConversation(ctx context.Context, req *conversation.UpdateConversationApiRequest) (*conversation.UpdateConversationApiResponse, error) {
	resp := new(conversation.UpdateConversationApiResponse)
	convID := req.GetConversationID()

	apiKeyInfo := ctxutil2.GetApiAuthFromCtx(ctx)
	userID := apiKeyInfo.UserID

	if userID == 0 {
		return resp, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "permission check failed"))
	}

	conversationDO, err := c.ConversationDomainSVC.GetByID(ctx, convID)
	if err != nil {
		return resp, err
	}
	if conversationDO == nil {
		return resp, errorx.New(errno.ErrConversationNotFound)
	}
	if conversationDO.CreatorID != userID {
		return resp, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "user not match"))
	}

	updateResult, err := c.ConversationDomainSVC.Update(ctx, &entity.UpdateMeta{
		ID:   convID,
		Name: req.GetName(),
	})
	if err != nil {
		return resp, err
	}
	resp.ConversationData = &conversation.ConversationData{
		Id:            updateResult.ID,
		LastSectionID: &updateResult.SectionID,
		CreatedAt:     updateResult.CreatedAt / 1000,
		Name:          ptr.Of(updateResult.Name),
	}
	return resp, nil
}

func parseMetaData(metaData map[string]string) string {
	if metaData == nil {
		return ""
	}
	j, err := json.Marshal(metaData)
	if err != nil {
		return ""
	}
	return string(j)
}

func parseExt(ext string) map[string]string {
	if ext == "" {
		return nil
	}
	var metaData map[string]string
	err := json.Unmarshal([]byte(ext), &metaData)
	if err != nil {
		return nil
	}
	return metaData
}

type OpenapiAgentRunApplication struct {
	UploadDomainSVC uploadService.UploadService
}

var ConversationOpenAPISVC = new(OpenapiAgentRunApplication)
