package entity

import (
	"fmt"
	"time"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type Message = model.Message
type Protocol = model.Protocol
type Role = model.Role
type ToolInfo = model.ToolInfo
type ToolChoice = model.ToolChoice
type ModelConfig = model.ModelConfig
type BizParam = model.BizParam

type ChatRequest struct {
	ModelConfig *ModelConfig
	Messages    []*Message
	Tools       []*ToolInfo
	BizParam    *BizParam
}

func (p *ChatRequest) GetMessages() (v []*Message) {
	return p.Messages
}

func (p *ChatRequest) GetTools() (v []*ToolInfo) {
	return p.Tools
}

func (p *ChatRequest) GetBizParam() (v *BizParam) {
	return p.BizParam
}

func (p *ChatRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ChatRequest(%+v)", *p)
}

type ChatResponse struct {
	Message *Message
}

func (p *ChatResponse) GetMessage() (v *Message) {
	return p.Message
}

func (p *ChatResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ChatResponse(%+v)", *p)
}

type ChatResponseStream interface {
	Send(*ChatResponse) error
}

type StreamRespParseResult struct {
	RespMsgs          []*Message
	LastRespMsg       *Message
	ReasoningDuration time.Duration
	FirstTokenLatency time.Duration
}
