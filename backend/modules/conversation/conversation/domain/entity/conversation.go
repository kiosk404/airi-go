package entity

import (
	"github.com/kiosk404/airi-go/backend/api/model/conversation/common"
	"github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/conversation/model"
)

type Conversation = model.Conversation

type CreateMeta struct {
	Name    string       `json:"name"`
	AgentID int64        `json:"agent_id"`
	UserID  int64        `json:"user_id"`
	Scene   common.Scene `json:"scene"`
	Ext     string       `json:"ext"`
}

type NewConversationCtxRequest struct {
	ID int64 `json:"id"`
}

type NewConversationCtxResponse struct {
	ID        int64 `json:"id"`
	SectionID int64 `json:"section_id"`
}

type GetCurrent = model.GetCurrent

type ListMeta struct {
	UserID  int64        `json:"user_id"`
	Scene   common.Scene `json:"scene"`
	AgentID int64        `json:"agent_id"`
	Limit   int          `json:"limit"`
	Page    int          `json:"page"`
}

type UpdateMeta struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
