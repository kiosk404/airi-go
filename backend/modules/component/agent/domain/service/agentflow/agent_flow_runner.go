package agentflow

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/agentrun/model"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/urltobase64url"
	"github.com/kiosk404/airi-go/backend/pkg/utils/safego"
)

type AgentState struct {
	Messages                 []*schema.Message
	UserInput                *schema.Message
	ReturnDirectlyToolCallID string
}

type AgentRequest struct {
	UserID  string
	Input   *schema.Message
	History []*schema.Message

	Identity *singleagent.AgentIdentity

	ResumeInfo   *singleagent.InterruptInfo // 中断恢复信息
	PreCallTools []*agentrun.ToolsRetriever // 预调用信息
	Variables    map[string]string          // 变量集合
}

type AgentRunner struct {
	runner            compose.Runnable[*AgentRequest, *schema.Message]
	requireCheckpoint bool

	returnDirectlyTools mapset.Set[string] // 直接返回的工具集
	containWfTool       bool               // 是否包含工作流
	modelInfo           *modelmgr.Model
}

func (r *AgentRunner) StreamExecute(ctx context.Context, req *AgentRequest) (sr *schema.StreamReader[*entity.AgentEvent], err error) {
	executeID := uuid.New()

	// 创建流式传输管道
	hdl, sr, sw := newReplyCallback(ctx, executeID.String(), r.returnDirectlyTools)
	var composeOpts []compose.Option
	//var pipeMsgOpt compose.Option
	//var workflowMsgSr *schema.StreamReader[*crossworkflow.WorkflowMessage]
	var workflowMsgCloser func()

	if r.containWfTool {
		// todo: 添加 workflow 相关
	}

	composeOpts = append(composeOpts, compose.WithCallbacks(hdl))
	_ = compose.RegisterSerializableType[*AgentState]("agent_state")
	if r.requireCheckpoint {
		defaultCheckPointID := executeID.String()
		if req.ResumeInfo != nil {
			resumeInfo := req.ResumeInfo
			if resumeInfo.InterruptType != singleagent.InterruptEventType_OauthPlugin {

			}
		}

		composeOpts = append(composeOpts, compose.WithCheckPointID(defaultCheckPointID))
	}

	safego.Go(ctx, func() {
		defer func() {
			if pe := recover(); pe != nil {
				logs.Error(pkg.ModelName, "[AgentRunner] StreamExecute recover, err: %v", pe)

				sw.Send(nil, errors.New("internal server error"))
			}
			if workflowMsgCloser != nil {
				workflowMsgCloser()
			}
			sw.Close()
		}()
		_, _ = r.runner.Stream(ctx, req, composeOpts...)
	})

	return sr, nil

}

func (r *AgentRunner) PreHandlerReq(ctx context.Context, req *AgentRequest) *AgentRequest {
	req.Input = r.preHandlerInput(req.Input)
	req.History = r.preHandlerHistory(req.History)
	logs.InfoX(pkg.ModelName, "[AgentRunner] PreHandlerReq, req: %v", conv.DebugJsonToStr(req))

	return req
}

func (r *AgentRunner) preHandlerInput(input *schema.Message) *schema.Message {
	var multiContent []schema.ChatMessagePart

	if len(input.MultiContent) == 0 {
		return input
	}

	unSupportMultiPart := make([]schema.ChatMessagePart, 0, len(input.MultiContent))

	for _, v := range input.MultiContent {
		switch v.Type {
		case schema.ChatMessagePartTypeImageURL:
			if !r.isSupportImage() {
				unSupportMultiPart = append(unSupportMultiPart, v)
			} else {
				v.ImageURL = transImageURLToBase64(v.ImageURL, r.enableLocalFileToLLMWithBase64())
				multiContent = append(multiContent, v)
			}
		case schema.ChatMessagePartTypeFileURL:
			if !r.isSupportFile() {
				unSupportMultiPart = append(unSupportMultiPart, v)
			} else {
				v.FileURL = transFileURLToBase64(v.FileURL, r.enableLocalFileToLLMWithBase64())
				multiContent = append(multiContent, v)
			}
		case schema.ChatMessagePartTypeAudioURL:
			if !r.isSupportAudio() {
				unSupportMultiPart = append(unSupportMultiPart, v)
			} else {
				v.AudioURL = transAudioURLToBase64(v.AudioURL, r.enableLocalFileToLLMWithBase64())
				multiContent = append(multiContent, v)
			}
		case schema.ChatMessagePartTypeVideoURL:
			if !r.isSupportVideo() {
				unSupportMultiPart = append(unSupportMultiPart, v)
			} else {
				v.VideoURL = transVideoURLToBase64(v.VideoURL, r.enableLocalFileToLLMWithBase64())
				multiContent = append(multiContent, v)
			}
		case schema.ChatMessagePartTypeText:
		default:
			multiContent = append(multiContent, v)
		}
	}

	for _, v := range input.MultiContent {
		if v.Type != schema.ChatMessagePartTypeText {
			continue
		}

		if r.isSupportMultiContent() {
			if len(multiContent) > 0 {
				v.Text = concatContentString(v.Text, unSupportMultiPart)
				multiContent = append(multiContent, v)
			} else {
				input.Content = concatContentString(v.Text, unSupportMultiPart)
			}
		} else {
			input.Content = concatContentString(v.Text, unSupportMultiPart)
		}

	}
	input.MultiContent = multiContent
	return input
}
func concatContentString(textContent string, unSupportTypeURL []schema.ChatMessagePart) string {
	if len(unSupportTypeURL) == 0 {
		return textContent
	}
	for _, v := range unSupportTypeURL {
		switch v.Type {
		case schema.ChatMessagePartTypeImageURL:
			textContent += "  this is a image:" + v.ImageURL.URL
		case schema.ChatMessagePartTypeFileURL:
			textContent += "  this is a file:" + v.FileURL.URL
		case schema.ChatMessagePartTypeAudioURL:
			textContent += "  this is a audio:" + v.AudioURL.URL
		case schema.ChatMessagePartTypeVideoURL:
			textContent += "  this is a video:" + v.VideoURL.URL
		default:
		}
	}
	return textContent
}

func (r *AgentRunner) preHandlerHistory(history []*schema.Message) []*schema.Message {
	var hm []*schema.Message
	for _, msg := range history {
		if msg.Role == schema.User {
			msg = r.preHandlerInput(msg)
		}
		hm = append(hm, msg)
	}
	return hm
}

func (r *AgentRunner) isSupportMultiContent() bool {
	return r.modelInfo.Capability != nil
}
func (r *AgentRunner) isSupportImage() bool { return r.modelInfo.Capability.GetImageUnderstanding() }
func (r *AgentRunner) isSupportFile() bool  { return r.modelInfo.Capability.GetAudioUnderstanding() }
func (r *AgentRunner) isSupportAudio() bool { return r.modelInfo.Capability.GetAudioUnderstanding() }
func (r *AgentRunner) isSupportVideo() bool { return r.modelInfo.Capability.GetVideoUnderstanding() }

func (r *AgentRunner) enableLocalFileToLLMWithBase64() bool {
	return r.modelInfo.EnableBase64URL
}

func transImageURLToBase64(imageUrl *schema.ChatMessageImageURL, enableBase64Url bool) *schema.ChatMessageImageURL {
	if !enableBase64Url {
		return imageUrl
	}
	fileData, err := urltobase64url.URLToBase64(imageUrl.URL)
	if err != nil {
		return imageUrl
	}
	imageUrl.URL = fileData.Base64Url
	imageUrl.MIMEType = fileData.MimeType
	return imageUrl
}

func transFileURLToBase64(fileUrl *schema.ChatMessageFileURL, enableBase64Url bool) *schema.ChatMessageFileURL {

	if !enableBase64Url {
		return fileUrl
	}
	fileData, err := urltobase64url.URLToBase64(fileUrl.URL)
	if err != nil {
		return fileUrl
	}
	fileUrl.URL = fileData.Base64Url
	fileUrl.MIMEType = fileData.MimeType
	return fileUrl
}

func transAudioURLToBase64(audioUrl *schema.ChatMessageAudioURL, enableBase64Url bool) *schema.ChatMessageAudioURL {

	if !enableBase64Url {
		return audioUrl
	}
	fileData, err := urltobase64url.URLToBase64(audioUrl.URL)
	if err != nil {
		return audioUrl
	}
	audioUrl.URL = fileData.Base64Url
	audioUrl.MIMEType = fileData.MimeType
	return audioUrl
}

func transVideoURLToBase64(videoUrl *schema.ChatMessageVideoURL, enableBase64Url bool) *schema.ChatMessageVideoURL {

	if !enableBase64Url {
		return videoUrl
	}
	fileData, err := urltobase64url.URLToBase64(videoUrl.URL)
	if err != nil {
		return videoUrl
	}
	videoUrl.URL = fileData.Base64Url
	videoUrl.MIMEType = fileData.MimeType
	return videoUrl
}
