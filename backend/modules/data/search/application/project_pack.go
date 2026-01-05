package application

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence"
	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence/common"
	"github.com/kiosk404/airi-go/backend/modules/data/search/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type projectInfo struct {
	iconURI string
	desc    string
}

type ProjectPacker interface {
	GetProjectInfo(ctx context.Context) (*projectInfo, error)
	GetPermissionInfo() *intelligence.IntelligencePermissionInfo
	GetPublishedInfo(ctx context.Context) *intelligence.IntelligencePublishInfo
	GetUserInfo(ctx context.Context, userID int64) *common.User
}

func NewPackProject(uid, projectID int64, tp common.IntelligenceType, s *SearchApplicationService) (ProjectPacker, error) {
	base := projectBase{SVC: s, projectID: projectID, iType: tp, uid: uid}

	switch tp {
	case common.IntelligenceType_Bot:
		return &agentPacker{projectBase: base}, nil
	case common.IntelligenceType_Project:
		return &appPacker{projectBase: base}, nil
	}

	return nil, fmt.Errorf("unsupported project_type: %d , project_id : %d", tp, projectID)
}

type projectBase struct {
	projectID int64 // agent_id or application_id
	uid       int64
	SVC       *SearchApplicationService
	iType     common.IntelligenceType
}

func (p *projectBase) GetPermissionInfo() *intelligence.IntelligencePermissionInfo {
	return &intelligence.IntelligencePermissionInfo{
		InCollaboration: false,
		CanDelete:       true,
		CanView:         true,
	}
}

func (p *projectBase) GetUserInfo(ctx context.Context, userID int64) *common.User {
	u, err := p.SVC.UserDomainSVC.GetUserInfo(ctx, userID)
	if err != nil {
		logs.ErrorX(pkg.ModelName, "[projectBase-GetUserInfo] failed to get user info, user_id: %d, err: %v", userID, err)
		return nil
	}

	return &common.User{
		UserID:         u.UserID,
		AvatarURL:      u.IconURL,
		UserUniqueName: u.UniqueName,
	}
}

type agentPacker struct {
	projectBase
}

func (a *agentPacker) GetProjectInfo(ctx context.Context) (*projectInfo, error) {
	agent, err := a.SVC.SingleAgentDomainSVC.GetSingleAgentDraft(ctx, a.projectID)
	if err != nil {
		return nil, err
	}

	if agent == nil {
		return nil, fmt.Errorf("agent info is nil")
	}
	return &projectInfo{
		iconURI: agent.IconURI,
		desc:    agent.Desc,
	}, nil
}

func (p *agentPacker) GetPublishedInfo(ctx context.Context) *intelligence.IntelligencePublishInfo {
	pubInfo, err := p.SVC.SingleAgentDomainSVC.GetPublishedInfo(ctx, p.projectID)
	if err != nil {
		logs.ErrorX(pkg.ModelName, "[agent-GetPublishedInfo]failed to get published info, agent_id: %d, err: %v", p.projectID, err)

		return nil
	}

	return &intelligence.IntelligencePublishInfo{
		PublishTime:  conv.Int64ToStr(pubInfo.LastPublishTimeMS / 1000),
		HasPublished: pubInfo.LastPublishTimeMS > 0,
	}
}

type appPacker struct {
	projectBase
}

func (a *appPacker) GetProjectInfo(ctx context.Context) (*projectInfo, error) {
	return &projectInfo{
		iconURI: "",
		desc:    "",
	}, nil
}

func (a *appPacker) GetPublishedInfo(ctx context.Context) *intelligence.IntelligencePublishInfo {
	return &intelligence.IntelligencePublishInfo{
		PublishTime:  "",
		HasPublished: false,
	}
}
