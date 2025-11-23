package convertor

import (
	druntime "github.com/kiosk404/airi-go/backend/api/model/llm/domain/runtime"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

func MessagesDTO2DO(dtos []*druntime.Message) (dos []*entity.Message) {
	return slices.Transform(dtos, func(dto *druntime.Message) *entity.Message {
		return MessageDTO2DO(dto)
	})
}

func MessageDTO2DO(dto *druntime.Message) (do *entity.Message) {
	if dto == nil {
		return nil
	}
	return &entity.Message{
		Role:              entity.Role(dto.GetRole()),
		Content:           dto.GetContent(),
		ReasoningContent:  dto.GetReasoningContent(),
		MultiModalContent: ChatMessagePartsDTO2DO(dto.GetMultimodalContents()),
		Name:              "",
		ToolCalls:         ToolCallsDTO2DO(dto.GetToolCalls()),
		ToolCallID:        dto.GetToolCallID(),
		ResponseMeta:      ResponseMetaDTO2DO(dto.GetResponseMeta()),
	}
}

func ResponseMetaDTO2DO(dto *druntime.ResponseMeta) (do *model.ResponseMeta) {
	if dto == nil {
		return nil
	}
	return &model.ResponseMeta{
		FinishReason: dto.GetFinishReason(),
		Usage:        TokenUsageDTO2DO(dto.GetUsage()),
	}
}

func TokenUsageDTO2DO(dto *druntime.TokenUsage) (do *model.TokenUsage) {
	if dto == nil {
		return nil
	}
	return &model.TokenUsage{
		PromptTokens:     int(dto.GetPromptTokens()),
		CompletionTokens: int(dto.GetCompletionTokens()),
		TotalTokens:      int(dto.GetTotalTokens()),
	}
}

func ToolCallsDTO2DO(dtos []*druntime.ToolCall) (dos []*model.ToolCall) {
	return slices.Map(dtos, func(dto *druntime.ToolCall, _ int) *model.ToolCall {
		return ToolCallDTO2DO(dto)
	})
}

func ToolCallDTO2DO(dto *druntime.ToolCall) (do *model.ToolCall) {
	if dto == nil {
		return nil
	}
	return &model.ToolCall{
		Index:    dto.Index,
		ID:       dto.GetID(),
		Type:     dto.GetType(),
		Function: FunctionCallDTO2DO(dto.GetFunctionCall()),
		Extra:    nil,
	}
}

func FunctionCallDTO2DO(dto *druntime.FunctionCall) (do *model.FunctionCall) {
	if dto == nil {
		return nil
	}
	return &model.FunctionCall{
		Name:      dto.GetName(),
		Arguments: dto.GetArguments(),
	}
}

func ChatMessagePartsDTO2DO(dtos []*druntime.ChatMessagePart) (dos []*model.ChatMessagePart) {
	return slices.Map(dtos, func(dto *druntime.ChatMessagePart, _ int) *model.ChatMessagePart {
		return ChatMessagePartDTO2DO(dto)
	})
}

func ChatMessagePartDTO2DO(dto *druntime.ChatMessagePart) (do *model.ChatMessagePart) {
	if dto == nil {
		return nil
	}
	return &model.ChatMessagePart{
		Type:     model.ChatMessagePartType(dto.GetType()),
		Text:     dto.GetText(),
		ImageURL: ChatMessageImageURLDTO2DO(dto.GetImageURL()),
	}
}

func ChatMessageImageURLDTO2DO(dto *druntime.ChatMessageImageURL) (do *model.ChatMessageImageURL) {
	if dto == nil {
		return nil
	}
	return &model.ChatMessageImageURL{
		URL:      dto.GetURL(),
		Detail:   model.ImageURLDetail(dto.GetDetail()),
		MIMEType: dto.GetMimeType(),
	}
}

func MessageDO2DTO(do *entity.Message) (dto *druntime.Message) {
	if do == nil {
		return nil
	}
	return &druntime.Message{
		Role:               druntime.Role(do.Role),
		Content:            ptr.Of(do.Content),
		ReasoningContent:   ptr.Of(do.ReasoningContent),
		MultimodalContents: ChatMessagePartsDO2DTO(do.MultiModalContent),
		ToolCalls:          ToolCallsDO2DTO(do.ToolCalls),
		ToolCallID:         ptr.Of(do.ToolCallID),
		ResponseMeta:       ResponseMetaDO2DTO(do.ResponseMeta),
	}
}

func ResponseMetaDO2DTO(do *model.ResponseMeta) (dto *druntime.ResponseMeta) {
	if do == nil {
		return nil
	}
	return &druntime.ResponseMeta{
		FinishReason: ptr.Of(do.FinishReason),
		Usage:        TokenUsageDO2DTO(do.Usage),
	}
}

func TokenUsageDO2DTO(do *model.TokenUsage) (dto *druntime.TokenUsage) {
	if do == nil {
		return nil
	}
	return &druntime.TokenUsage{
		PromptTokens:     ptr.Of(int64(do.PromptTokens)),
		CompletionTokens: ptr.Of(int64(do.CompletionTokens)),
		TotalTokens:      ptr.Of(int64(do.TotalTokens)),
	}
}

func ToolCallsDO2DTO(dos []*model.ToolCall) (dtos []*druntime.ToolCall) {
	return slices.Map(dos, func(do *model.ToolCall, _ int) *druntime.ToolCall {
		return ToolCallDO2DTO(do)
	})
}

func ToolCallDO2DTO(do *model.ToolCall) (dto *druntime.ToolCall) {
	if do == nil {
		return nil
	}
	return &druntime.ToolCall{
		Index:        do.Index,
		ID:           ptr.Of(do.ID),
		Type:         ptr.Of(do.Type),
		FunctionCall: FunctionCallDO2DTO(do.Function),
	}
}

func FunctionCallDO2DTO(do *model.FunctionCall) (dto *druntime.FunctionCall) {
	if do == nil {
		return nil
	}
	return &druntime.FunctionCall{
		Name:      ptr.Of(do.Name),
		Arguments: ptr.Of(do.Arguments),
	}
}

func ChatMessagePartsDO2DTO(dos []*model.ChatMessagePart) (dtos []*druntime.ChatMessagePart) {
	return slices.Map(dos, func(do *model.ChatMessagePart, _ int) *druntime.ChatMessagePart {
		return ChatMessagePartDO2DTO(do)
	})
}

func ChatMessagePartDO2DTO(do *model.ChatMessagePart) (dto *druntime.ChatMessagePart) {
	if do == nil {
		return nil
	}
	return &druntime.ChatMessagePart{
		Type:     ptr.Of(druntime.ChatMessagePartType(do.Type)),
		Text:     ptr.Of(do.Text),
		ImageURL: ChatMessageImageURLDO2DTO(do.ImageURL),
	}
}

func ChatMessageImageURLDO2DTO(do *model.ChatMessageImageURL) (dto *druntime.ChatMessageImageURL) {
	if do == nil {
		return nil
	}
	return &druntime.ChatMessageImageURL{
		URL:      ptr.Of(do.URL),
		Detail:   ptr.Of(druntime.ImageURLDetail(do.Detail)),
		MimeType: ptr.Of(do.MIMEType),
	}
}
