package message

import (
	"context"

	"github.com/cloudwego/eino/schema"
	model "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
)

//go:generate  mockgen -destination messagemock/message_mock.go --package messagemock -source message.go
type Message interface {
	GetByRunIDs(ctx context.Context, conversationID int64, runIDs []int64) ([]*model.Message, error)
	PreCreate(ctx context.Context, msg *model.Message) (*model.Message, error)
	Create(ctx context.Context, msg *model.Message) (*model.Message, error)
	BatchCreate(ctx context.Context, msg []*model.Message) ([]*model.Message, error)
	List(ctx context.Context, meta *entity.ListMeta) (*entity.ListResult, error)
	ListWithoutPair(ctx context.Context, req *entity.ListMeta) (*entity.ListResult, error)
	Edit(ctx context.Context, msg *model.Message) (*model.Message, error)
	Delete(ctx context.Context, req *entity.DeleteMeta) error
	GetMessageByID(ctx context.Context, id int64) (*entity.Message, error)
	MessageList(ctx context.Context, req *MessageListRequest) (*MessageListResponse, error)
	GetLatestRunIDs(ctx context.Context, req *GetLatestRunIDsRequest) ([]int64, error)
	GetMessagesByRunIDs(ctx context.Context, req *GetMessagesByRunIDsRequest) (*GetMessagesByRunIDsResponse, error)
}

var defaultSVC Message

type MessageMeta = model.Message

func DefaultSVC() Message {
	return defaultSVC
}

func SetDefaultSVC(c Message) {
	defaultSVC = c
}

type MessageListRequest struct {
	ConversationID int64
	Limit          int64
	BeforeID       *string
	AfterID        *string
	UserID         int64
	BizID          int64
	OrderBy        *string
}

type MessageListResponse struct {
	Messages []*WfMessage
	FirstID  string
	LastID   string
	HasMore  bool
}

type Content struct {
	Type model.InputType `json:"type"`
	Text *string         `json:"text,omitempty"`
	Uri  *string         `json:"uri,omitempty"`
	Url  *string         `json:"url,omitempty"`
}

type WfMessage struct {
	ID           int64
	Role         schema.RoleType `json:"role"` // user or assistant
	MultiContent []*Content      `json:"multi_content"`
	Text         *string         `json:"text,omitempty"`
	ContentType  string          `json:"content_type"`
	SectionID    int64           `json:"section_id"`
}

type GetLatestRunIDsRequest struct {
	ConversationID int64
	UserID         int64
	BizID          int64
	Rounds         int64
	SectionID      int64
	InitRunID      *int64
}

type GetMessagesByRunIDsRequest struct {
	ConversationID int64
	RunIDs         []int64
}

type GetMessagesByRunIDsResponse struct {
	Messages       []*WfMessage
	SchemaMessages []*schema.Message
}
