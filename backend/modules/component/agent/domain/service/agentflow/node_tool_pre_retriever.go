package agentflow

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

type toolPreCallConf struct{}

func newPreToolRetriever(conf *toolPreCallConf) *toolPreCallConf {
	return &toolPreCallConf{}
}

// 在 Agent 执行主逻辑之前，预先执行一些指定的工具，并将工具的结果作为上下文信息给 Agent，帮助 Agent 做出判断
// PreCallTools:
// - 工具1：查询用户信息（是否VIP）
// - 工具2：查询库存
// 执行过程：
// 1. 执行工具1（查询用户信息）
// 2. 构建消息对（1. Assistant: ToolCall：{"用户是什么身份"}, 2. Tool: Content: "用户是VIP"）
// 3. 执行工具2（查询库存）
// 4. 构建消息对（3. Assistant: ToolCall: {"库存还有吗？"}, 4. Tool: Content: "库存充足，可以立即发货"）
// 返回所有的消息 [消息1，消息2，消息3，消息4]
func (pr *toolPreCallConf) toolPreRetrieve(ctx context.Context, ar *AgentRequest) ([]*schema.Message, error) {
	var tms []*schema.Message

	return tms, nil
}
