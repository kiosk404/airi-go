package entity

import (
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/pkg/errors"
)

func FromDOMessages(dos []*Message) []*schema.Message {
	return slices.Map(dos, func(do *Message, _ int) *schema.Message {
		return FromDOMessage(do)
	})
}

func FromDOMessage(do *Message) *schema.Message {
	if do == nil {
		return nil
	}
	return &schema.Message{
		Role:         schema.RoleType(do.Role),
		Content:      do.Content,
		MultiContent: FromDOChatMsgParts(do.MultiModalContent),
		Name:         do.Name,
		ToolCalls:    FromDOToolCalls(do.ToolCalls),
		ToolCallID:   do.ToolCallID,
		ResponseMeta: FromDOResponseMeta(do.ResponseMeta),
		// Extra:        nil,
	}
}

func FromDOChatMsgParts(ps []*model.ChatMessagePart) []schema.ChatMessagePart {
	return slices.Map(ps, func(p *model.ChatMessagePart, _ int) schema.ChatMessagePart {
		return FromDOChatMsgPart(p)
	})
}

func FromDOChatMsgPart(p *model.ChatMessagePart) schema.ChatMessagePart {
	if p == nil {
		return schema.ChatMessagePart{}
	}
	return schema.ChatMessagePart{
		Type:     schema.ChatMessagePartType(p.Type),
		Text:     p.Text,
		ImageURL: FromDOImageURL(p.ImageURL),
	}
}

func FromDOImageURL(p *model.ChatMessageImageURL) *schema.ChatMessageImageURL {
	if p == nil {
		return nil
	}
	return &schema.ChatMessageImageURL{
		URL:      p.URL,
		Detail:   schema.ImageURLDetail(p.Detail),
		MIMEType: p.MIMEType,
	}
}

func FromDOToolCalls(ts []*model.ToolCall) []schema.ToolCall {
	return slices.Map(ts, func(t *model.ToolCall, _ int) schema.ToolCall {
		return FromDOToolCall(t)
	})
}

func FromDOToolCall(t *model.ToolCall) schema.ToolCall {
	if t == nil {
		return schema.ToolCall{}
	}
	return schema.ToolCall{
		Index: ptr.PtrConvertMap(t.Index, func(f int64) int {
			return int(f)
		}),
		ID:       t.ID,
		Type:     t.Type,
		Function: FromDOFunctionCall(t.Function),
	}
}

func FromDOFunctionCall(f *model.FunctionCall) schema.FunctionCall {
	if f == nil {
		return schema.FunctionCall{}
	}
	return schema.FunctionCall{
		Name:      f.Name,
		Arguments: f.Arguments,
	}
}

func FromDOResponseMeta(rm *model.ResponseMeta) *schema.ResponseMeta {
	if rm == nil {
		return nil
	}
	return &schema.ResponseMeta{
		FinishReason: rm.FinishReason,
		Usage:        FromDOTokenUsage(rm.Usage),
	}
}

func FromDOTokenUsage(tu *model.TokenUsage) *schema.TokenUsage {
	if tu == nil {
		return nil
	}
	return &schema.TokenUsage{
		PromptTokens:     tu.PromptTokens,
		CompletionTokens: tu.CompletionTokens,
		TotalTokens:      tu.TotalTokens,
	}
}

func FromDOOptions(options *model.Options) ([]einoModel.Option, error) {
	var res []einoModel.Option
	if options.Temperature != nil {
		res = append(res, einoModel.WithTemperature(*options.Temperature))
	}
	if options.MaxTokens != nil {
		res = append(res, einoModel.WithMaxTokens(*options.MaxTokens))
	}
	if options.Model != nil {
		res = append(res, einoModel.WithModel(*options.Model))
	}
	if options.TopP != nil {
		res = append(res, einoModel.WithTopP(*options.TopP))
	}
	if options.Stop != nil {
		res = append(res, einoModel.WithStop(options.Stop))
	}
	if options.ToolChoice != nil {
		res = append(res, einoModel.WithToolChoice(FromDOToolChoice(*options.ToolChoice)))
	}
	return res, nil
}

func FromDOToolChoice(do model.ToolChoice) (einoToolChoice schema.ToolChoice) {
	switch do {
	case model.ToolChoiceNone:
		return schema.ToolChoiceForbidden
	case model.ToolChoiceAuto:
		return schema.ToolChoiceAllowed
	case model.ToolChoiceRequired:
		return schema.ToolChoiceForced
	default:
		return
	}
}

func FromDOTools(dos []*model.ToolInfo) ([]*schema.ToolInfo, error) {
	if len(dos) == 0 {
		return nil, nil
	}
	res := make([]*schema.ToolInfo, 0)
	for _, do := range dos {
		einoTool, err := FromDOTool(do)
		if err != nil {
			return nil, err
		}
		res = append(res, einoTool)
	}
	return res, nil
}

func FromDOTool(do *model.ToolInfo) (*schema.ToolInfo, error) {
	if do == nil {
		return nil, nil
	}
	if do.ToolDefType != model.ToolDefTypeOpenAPIV3 {
		return nil, errors.Errorf("unsupport tool def type:%s", do.ToolDefType)
	}
	var openApiV3Schema openapi3.Schema
	if err := sonic.UnmarshalString(do.Def, &openApiV3Schema); err != nil {
		return nil, errors.Errorf("[fromDOTool] unmarshal tool def failed, err:%s", err.Error())
	}
	return &schema.ToolInfo{
		Name:        do.Name,
		Desc:        do.Desc,
		ParamsOneOf: schema.NewParamsOneOfByOpenAPIV3(&openApiV3Schema),
	}, nil
}

func ToDOMessages(msgs []*schema.Message) ([]*Message, error) {
	if len(msgs) == 0 {
		return nil, nil
	}
	res := make([]*Message, len(msgs))
	for i, msg := range msgs {
		do, err := ToDOMessage(msg)
		if err != nil {
			return nil, err
		}
		res[i] = do
	}
	return res, nil
}

func ToDOMessage(msg *schema.Message) (*Message, error) {
	if msg == nil {
		return nil, nil
	}
	return &Message{
		Role:              Role(msg.Role),
		Content:           msg.Content,
		ReasoningContent:  GetReasoningContent(msg),
		MultiModalContent: ToDOMultiContents(msg.MultiContent),
		Name:              msg.Name,
		ToolCalls:         ToDOToolCalls(msg.ToolCalls),
		ToolCallID:        msg.ToolCallID,
		ResponseMeta:      ToDORespMeta(msg.ResponseMeta),
	}, nil
}

func GetReasoningContent(msg *schema.Message) string {
	rc, ok := deepseek.GetReasoningContent(msg)
	if ok {
		return rc
	}
	return ""
}

func ToDOToolCalls(tcs []schema.ToolCall) []*model.ToolCall {
	return slices.Map(tcs, func(tc schema.ToolCall, _ int) *model.ToolCall {
		return ToDOToolCall(tc)
	})
}

func ToDOToolCall(tc schema.ToolCall) *model.ToolCall {
	return &model.ToolCall{
		Index: ptr.PtrConvertMap(tc.Index, func(f int) int64 {
			return int64(f)
		}),
		ID:       tc.ID,
		Type:     tc.Type,
		Function: ToDOFunctionCall(tc.Function),
		Extra:    tc.Extra,
	}
}

func ToDOFunctionCall(f schema.FunctionCall) *model.FunctionCall {
	return &model.FunctionCall{
		Name:      f.Name,
		Arguments: f.Arguments,
	}
}

func ToDOMultiContents(cms []schema.ChatMessagePart) []*model.ChatMessagePart {
	return slices.Map(cms, func(cm schema.ChatMessagePart, _ int) *model.ChatMessagePart {
		return ToDOMultiContent(cm)
	})
}

func ToDOMultiContent(cm schema.ChatMessagePart) *model.ChatMessagePart {
	return &model.ChatMessagePart{
		Type:     model.ChatMessagePartType(cm.Type),
		Text:     cm.Text,
		ImageURL: ToDOImageURL(cm.ImageURL),
	}
}

func ToDOImageURL(cm *schema.ChatMessageImageURL) *model.ChatMessageImageURL {
	if cm == nil {
		return nil
	}
	return &model.ChatMessageImageURL{
		URL:      cm.URL,
		Detail:   model.ImageURLDetail(cm.Detail),
		MIMEType: cm.MIMEType,
	}
}

func ToDORespMeta(rm *schema.ResponseMeta) *model.ResponseMeta {
	if rm == nil {
		return nil
	}
	return &model.ResponseMeta{
		FinishReason: rm.FinishReason,
		Usage:        ToDOTokenUsage(rm.Usage),
	}
}

func ToDOTokenUsage(tu *schema.TokenUsage) *model.TokenUsage {
	if tu == nil {
		return nil
	}
	return &model.TokenUsage{
		PromptTokens:     tu.PromptTokens,
		CompletionTokens: tu.CompletionTokens,
		TotalTokens:      tu.TotalTokens,
	}
}
