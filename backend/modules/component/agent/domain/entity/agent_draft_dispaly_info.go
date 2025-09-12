package entity

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
)

type AgentDraftDisplayInfo struct {
	AgentID     int64
	DisplayInfo *developer_api.DraftBotDisplayInfoData
}
