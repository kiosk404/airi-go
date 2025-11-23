package agentflow

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
)

type variableConf struct {
	Agent  *entity.SingleAgent
	UserID string
}

func loadAgentVariables(ctx context.Context, vc *variableConf) (map[string]string, error) {
	vbs := make(map[string]string)

	// todo:// 暂不支持变量

	return vbs, nil
}

func newAgentVariableTools(ctx context.Context, v *variableConf) ([]tool.InvokableTool, error) {
	tools := make([]tool.InvokableTool, 0, 1)
	a := &avTool{
		Agent:  v.Agent,
		UserID: v.UserID,
	}

	desc := `
## Skills Conditions
1. When the user's intention is to set a variable and the user provides the variable to be set, call the tool.
2. If the user wants to set a variable but does not provide the variable, do not call the tool.
3. If the user's intention is not to set a variable, do not call the tool.

## Constraints
- Only make decisions regarding tool invocation based on the user's intention and input related to variable setting.
- Do not call the tool in any other situation not meeting the above conditions.
`
	at, err := utils.InferTool("setKeywordMemory", desc, a.Invoke)
	if err != nil {
		return nil, err
	}
	tools = append(tools, at)
	return tools, nil
}

type avTool struct {
	Agent  *entity.SingleAgent
	UserID string
}

type KVMeta struct {
	Keyword string `json:"keyword" jsonschema:"required,description=the keyword of memory variable"`
	Value   string `json:"value" jsonschema:"required,description=the value of memory variable"`
}
type KVMemoryVariable struct {
	Data []*KVMeta `json:"data"`
}

func (a *avTool) Invoke(ctx context.Context, v *KVMemoryVariable) (string, error) {

	return "success", nil
}
