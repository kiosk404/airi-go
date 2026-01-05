package application

import (
	"context"
	"sync"

	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence"
	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence/common"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	agentEntity "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/data/search/pkg"
	"github.com/kiosk404/airi-go/backend/modules/data/search/pkg/errno"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/pkg/consts"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/taskgroup"
)

var projectType2iconURI = map[common.IntelligenceType]string{
	common.IntelligenceType_Bot:     consts.DefaultAgentIcon,
	common.IntelligenceType_Project: consts.DefaultAppIcon,
}

// GetDraftIntelligenceList 基于 MySQL 查询草稿智能体列表
// 重构说明: 原实现基于 ES 查询，现改为直接查询 MySQL，后续可再考虑 ES 优化
func (s *SearchApplicationService) GetDraftIntelligenceList(ctx context.Context, req *intelligence.GetDraftIntelligenceListRequest) (
	resp *intelligence.GetDraftIntelligenceListResponse, err error,
) {
	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		logs.ErrorX(pkg.ModelName, "[GetDraftIntelligenceList] userID is nil, session is required")
		return nil, errorx.New(errno.ErrSearchPermissionCode, errorx.KV("msg", "session is required"))
	}

	logs.InfoX(pkg.ModelName, "[GetDraftIntelligenceList] userID: %d, req.Types: %v, req.Size: %d", *userID, req.GetTypes(), req.GetSize())

	// 计算分页参数
	pageSize := int(req.GetSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	page := 1
	// cursor_id 暂不支持，使用简单分页
	// TODO: 后续可以基于 cursor_id 实现游标分页

	// 直接从 MySQL 查询 Agent 列表
	agents, total, err := s.SingleAgentDomainSVC.ListAgentDraftByCreator(ctx, *userID, page, pageSize)
	if err != nil {
		logs.ErrorX(pkg.ModelName, "[GetDraftIntelligenceList] ListAgentDraftByCreator failed: %v", err)
		return nil, errorx.Wrapf(err, "ListAgentDraftByCreator failed")
	}

	logs.InfoX(pkg.ModelName, "[GetDraftIntelligenceList] query result: userID=%d, agents count=%d, total=%d", *userID, len(agents), total)

	if len(agents) == 0 {
		return &intelligence.GetDraftIntelligenceListResponse{
			Code: 0,
			Data: &intelligence.DraftIntelligenceListData{
				Intelligences: make([]*intelligence.IntelligenceData, 0),
				Total:         0,
				HasMore:       false,
				NextCursorID:  "",
			},
		}, nil
	}

	// 并发打包 Intelligence 数据
	intelligenceDataList := make([]*intelligence.IntelligenceData, len(agents))
	tasks := taskgroup.NewUninterruptibleTaskGroup(ctx, len(agents))
	lock := sync.Mutex{}

	for idx, agent := range agents {
		index := idx
		agentData := agent
		tasks.Go(func() error {
			info, err := s.packAgentToIntelligenceData(ctx, agentData)
			if err != nil {
				logs.ErrorX(pkg.ModelName, "[packAgentToIntelligenceData] failed id %v, name %s, err: %v",
					agentData.AgentID, agentData.Name, err)
				return nil
			}

			lock.Lock()
			defer lock.Unlock()
			intelligenceDataList[index] = info
			return nil
		})
	}

	if err = tasks.Wait(); err != nil {
		return nil, err
	}

	// 过滤掉 nil 数据
	filterDataList := make([]*intelligence.IntelligenceData, 0, len(intelligenceDataList))
	for _, data := range intelligenceDataList {
		if data != nil {
			filterDataList = append(filterDataList, data)
		}
	}

	// 计算是否有更多数据
	hasMore := int64(len(agents)) < total

	return &intelligence.GetDraftIntelligenceListResponse{
		Code: 0,
		Data: &intelligence.DraftIntelligenceListData{
			Intelligences: filterDataList,
			Total:         int32(total),
			HasMore:       hasMore,
			NextCursorID:  "", // 简单分页暂不支持 cursor
		},
	}, nil
}

// packAgentToIntelligenceData 将 SingleAgent 实体转换为 IntelligenceData
func (s *SearchApplicationService) packAgentToIntelligenceData(ctx context.Context, agent *agentEntity.SingleAgent) (*intelligence.IntelligenceData, error) {
	if agent == nil {
		return nil, nil
	}

	uid := ctxutil.MustGetUIDFromCtx(ctx)

	// 构建基础信息
	intelligenceData := &intelligence.IntelligenceData{
		Type: common.IntelligenceType_Bot,
		BasicInfo: &common.IntelligenceBasicInfo{
			ID:          agent.AgentID,
			Name:        agent.Name,
			Description: agent.Desc,
			IconURI:     agent.IconURI,
			IconURL:     s.getProjectIconURL(ctx, agent.IconURI, common.IntelligenceType_Bot),
			OwnerID:     agent.CreatorID,
			Status:      common.IntelligenceStatus_Using,
			CreateTime:  agent.CreatedAt,
			UpdateTime:  agent.UpdatedAt,
		},
	}

	// 获取权限信息
	intelligenceData.PermissionInfo = &intelligence.IntelligencePermissionInfo{
		InCollaboration: false,
		CanDelete:       true,
		CanView:         true,
	}

	// 获取发布信息
	pubInfo, err := s.SingleAgentDomainSVC.GetPublishedInfo(ctx, agent.AgentID)
	if err != nil {
		logs.WarnX(pkg.ModelName, "[packAgentToIntelligenceData] GetPublishedInfo failed, agent_id: %d, err: %v", agent.AgentID, err)
		intelligenceData.PublishInfo = &intelligence.IntelligencePublishInfo{
			HasPublished: false,
		}
	} else if pubInfo != nil && pubInfo.LastPublishTimeMS > 0 {
		intelligenceData.PublishInfo = &intelligence.IntelligencePublishInfo{
			PublishTime:  conv.Int64ToStr(pubInfo.LastPublishTimeMS / 1000),
			HasPublished: true,
		}
		intelligenceData.BasicInfo.PublishTime = pubInfo.LastPublishTimeMS / 1000
	} else {
		intelligenceData.PublishInfo = &intelligence.IntelligencePublishInfo{
			HasPublished: false,
		}
	}

	// 获取用户信息
	if s.UserDomainSVC != nil {
		userInfo, err := s.UserDomainSVC.GetUserInfo(ctx, agent.CreatorID)
		if err != nil {
			logs.WarnX(pkg.ModelName, "[packAgentToIntelligenceData] GetUserInfo failed, user_id: %d, err: %v", agent.CreatorID, err)
		} else if userInfo != nil {
			intelligenceData.OwnerInfo = &common.User{
				UserID:         userInfo.UserID,
				AvatarURL:      userInfo.IconURL,
				UserUniqueName: userInfo.UniqueName,
			}
		}
	}

	// 其他信息
	intelligenceData.LatestAuditInfo = &common.AuditInfo{}
	intelligenceData.FavoriteInfo = &intelligence.FavoriteInfo{
		IsFav:   false,
		FavTime: "",
	}
	intelligenceData.OtherInfo = &intelligence.OtherInfo{
		BotMode:          intelligence.BotMode_SingleMode,
		RecentlyOpenTime: "",
	}

	_ = uid // 保留 uid 以备后续使用

	return intelligenceData, nil
}

func (s *SearchApplicationService) getProjectIconURL(ctx context.Context, uri string, tp common.IntelligenceType) string {
	if uri == "" {
		return s.getProjectDefaultIconURL(ctx, tp)
	}

	url := s.getURL(ctx, uri)
	if url != "" {
		return url
	}

	return s.getProjectDefaultIconURL(ctx, tp)
}

func (s *SearchApplicationService) getProjectDefaultIconURL(ctx context.Context, tp common.IntelligenceType) string {
	iconURL, ok := projectType2iconURI[tp]
	if !ok {
		logs.WarnX(pkg.ModelName, "[getProjectDefaultIconURL] don't have type: %d  default icon", tp)

		return ""
	}

	return s.getURL(ctx, iconURL)
}
