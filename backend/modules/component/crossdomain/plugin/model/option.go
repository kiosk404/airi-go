package model

import (
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/consts"
)

type ExecuteToolOption struct {
	ProjectInfo *ProjectInfo

	AutoGenRespSchema bool

	ToolVersion                string
	Operation                  *Openapi3Operation
	InvalidRespProcessStrategy consts.InvalidResponseProcessStrategy

	ConversationID int64
}

type ExecuteToolOpt func(o *ExecuteToolOption)

type ProjectInfo struct {
	ProjectID      int64              // agentID or appID
	ProjectVersion *string            // if version si nil, use latest version
	ProjectType    consts.ProjectType // agent or app
}

func WithProjectInfo(info *ProjectInfo) ExecuteToolOpt {
	return func(o *ExecuteToolOption) {
		o.ProjectInfo = info
	}
}

func WithToolVersion(version string) ExecuteToolOpt {
	return func(o *ExecuteToolOption) {
		o.ToolVersion = version
	}
}

func WithOpenapiOperation(op *Openapi3Operation) ExecuteToolOpt {
	return func(o *ExecuteToolOption) {
		o.Operation = op
	}
}

func WithInvalidRespProcessStrategy(strategy consts.InvalidResponseProcessStrategy) ExecuteToolOpt {
	return func(o *ExecuteToolOption) {
		o.InvalidRespProcessStrategy = strategy
	}
}

func WithAutoGenRespSchema() ExecuteToolOpt {
	return func(o *ExecuteToolOption) {
		o.AutoGenRespSchema = true
	}
}

func WithPluginHTTPHeader(conversationID int64) ExecuteToolOpt {
	return func(o *ExecuteToolOption) {
		o.ConversationID = conversationID
	}
}
