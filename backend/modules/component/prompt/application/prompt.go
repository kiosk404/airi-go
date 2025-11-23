package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/playground"
	"github.com/kiosk404/airi-go/backend/api/model/resource/common"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/entity"
	prompt "github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/pkg"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/pkg/errno"
	searchEntity "github.com/kiosk404/airi-go/backend/modules/data/search/domain/entity"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type PromptApplicationService struct {
	DomainSVC prompt.Prompt
	eventbus  search.ResourceEventBus
}

var PromptSVC = &PromptApplicationService{}

func (p *PromptApplicationService) UpsertPromptResource(ctx context.Context, req *playground.UpsertPromptResourceRequest) (resp *playground.UpsertPromptResourceResponse, err error) {
	session := ctxutil.GetUserSessionFromCtx(ctx)
	if session == nil {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no session data provided"))
	}

	promptID := req.Prompt.GetID()
	if promptID == 0 {
		// create a new prompt resource
		resp, err = p.createPromptResource(ctx, req)
		if err != nil {
			return nil, err
		}

		pErr := p.eventbus.PublishResources(ctx, &searchEntity.ResourceDomainEvent{
			OpType: searchEntity.Created,
			Resource: &searchEntity.ResourceDocument{
				ResType:       common.ResType_Prompt,
				ResID:         resp.Data.ID,
				Name:          req.Prompt.Name,
				OwnerID:       &session.UserID,
				PublishStatus: ptr.Of(common.PublishStatus_Published),
			},
		})
		if pErr != nil {
			logs.ErrorX(pkg.ModelName, "publish resource event failed: %v", pErr)
		}

		return resp, nil
	}

	// update an existing prompt resource
	resp, err = p.updatePromptResource(ctx, req)
	if err != nil {
		return nil, err
	}

	pErr := p.eventbus.PublishResources(ctx, &searchEntity.ResourceDomainEvent{
		OpType: searchEntity.Updated,
		Resource: &searchEntity.ResourceDocument{
			ResType: common.ResType_Prompt,
			ResID:   resp.Data.ID,
			Name:    req.Prompt.Name,
		},
	})
	if pErr != nil {
		logs.ErrorX(pkg.ModelName, "publish resource event failed: %v", pErr)
	}

	return resp, nil
}

func (p *PromptApplicationService) GetPromptResourceInfo(ctx context.Context, req *playground.GetPromptResourceInfoRequest) (
	resp *playground.GetPromptResourceInfoResponse, err error,
) {

	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no session data provided"))
	}

	promptInfo, err := p.DomainSVC.GetPromptResource(ctx, req.GetPromptResourceID())
	if err != nil {
		return nil, err
	}
	if promptInfo.CreatorID != *uid {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no permission"))
	}

	return &playground.GetPromptResourceInfoResponse{
		Data: promptInfoDo2To(promptInfo),
		Code: 0,
	}, nil
}

func (p *PromptApplicationService) GetOfficialPromptResourceList(ctx context.Context, c *playground.GetOfficialPromptResourceListRequest) (
	*playground.GetOfficialPromptResourceListResponse, error,
) {
	session := ctxutil.GetUserSessionFromCtx(ctx)
	if session == nil {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no session data provided"))
	}

	promptList, err := p.DomainSVC.ListOfficialPromptResource(ctx, c.GetKeyword())
	if err != nil {
		return nil, err
	}

	return &playground.GetOfficialPromptResourceListResponse{
		PromptResourceList: slices.Transform(promptList, func(p *entity.PromptResource) *playground.PromptResource {
			return promptInfoDo2To(p)
		}),
		Code: 0,
	}, nil
}

func (p *PromptApplicationService) DeletePromptResource(ctx context.Context, req *playground.DeletePromptResourceRequest) (resp *playground.DeletePromptResourceResponse, err error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no session data provided"))
	}

	promptInfo, err := p.DomainSVC.GetPromptResource(ctx, req.GetPromptResourceID())
	if err != nil {
		return nil, err
	}

	if promptInfo.CreatorID != *uid {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no permission"))
	}

	err = p.DomainSVC.DeletePromptResource(ctx, req.GetPromptResourceID())
	if err != nil {
		return nil, err
	}

	pErr := p.eventbus.PublishResources(ctx, &searchEntity.ResourceDomainEvent{
		OpType: searchEntity.Deleted,
		Resource: &searchEntity.ResourceDocument{
			ResType: common.ResType_Prompt,
			ResID:   req.GetPromptResourceID(),
		},
	})
	if pErr != nil {
		logs.ErrorX(pkg.ModelName, "publish resource event failed: %v", pErr)
	}

	return &playground.DeletePromptResourceResponse{
		Code: 0,
	}, nil
}

func (p *PromptApplicationService) createPromptResource(ctx context.Context, req *playground.UpsertPromptResourceRequest) (resp *playground.UpsertPromptResourceResponse, err error) {
	do := p.toPromptResourceDO(req.Prompt)
	uid := ctxutil.GetUIDFromCtx(ctx)

	do.CreatorID = *uid

	promptID, err := p.DomainSVC.CreatePromptResource(ctx, do)
	if err != nil {
		return nil, err
	}

	return &playground.UpsertPromptResourceResponse{
		Data: &playground.ShowPromptResource{
			ID: promptID,
		},
		Code: 0,
	}, nil
}

func (p *PromptApplicationService) updatePromptResource(ctx context.Context, req *playground.UpsertPromptResourceRequest) (resp *playground.UpsertPromptResourceResponse, err error) {
	promptID := req.Prompt.GetID()

	promptResource, err := p.DomainSVC.GetPromptResource(ctx, promptID)
	if err != nil {
		return nil, err
	}

	logs.InfoX(pkg.ModelName, "promptResource.CreatorID : %v", promptResource.CreatorID)
	uid := ctxutil.GetUIDFromCtx(ctx)

	if promptResource.CreatorID != *uid {
		return nil, errorx.New(errno.ErrPromptPermissionCode, errorx.KV("msg", "no permission"))
	}

	err = p.DomainSVC.UpdatePromptResource(ctx, promptID, req.Prompt.Name, req.Prompt.Description, req.Prompt.PromptText)
	if err != nil {
		return nil, err
	}

	return &playground.UpsertPromptResourceResponse{
		Data: &playground.ShowPromptResource{
			ID: promptID,
		},
		Code: 0,
	}, nil
}

func (p *PromptApplicationService) toPromptResourceDO(m *playground.PromptResource) *entity.PromptResource {
	e := entity.PromptResource{}
	e.ID = m.GetID()
	e.PromptText = m.GetPromptText()
	e.Name = m.GetName()
	e.Description = m.GetDescription()

	return &e
}

func promptInfoDo2To(p *entity.PromptResource) *playground.PromptResource {
	return &playground.PromptResource{
		ID:          ptr.Of(p.ID),
		Name:        ptr.Of(p.Name),
		Description: ptr.Of(p.Description),
		PromptText:  ptr.Of(p.PromptText),
	}
}
