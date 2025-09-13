package entity

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
)

type PublishConnectorData struct {
	PublishConnectorList  []*developer_api.PublishConnectorInfo
	SubmitBotMarketOption *developer_api.SubmitBotMarketOption
	LastSubmitConfig      *developer_api.SubmitBotMarketConfig
	PublishTips           *developer_api.PublishTips
}

type PublishInfo struct {
	AgentID           int64
	LastPublishTimeMS int64
}
