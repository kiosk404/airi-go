package agentflow

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
)

const (
	placeholderOfUserInput   = "_user_input"
	placeholderOfChatHistory = "_chat_history"
)

type promptVariables struct {
	Agent *entity.SingleAgent
	avs   map[string]string
}

func (p *promptVariables) AssemblePromptVariables(ctx context.Context, req *AgentRequest) (variables map[string]any, err error) {
	variables = make(map[string]any)

	variables[placeholderOfTime] = time.Now().Format("Monday 2006/01/02 15:04:05 -07")
	variables[placeholderOfAgentName] = p.Agent.Name

	if req.Input != nil {
		variables[placeholderOfUserInput] = []*schema.Message{req.Input}
	}

	// Handling conversation history
	if len(req.History) > 0 {
		// Add chat history to variable
		variables[placeholderOfChatHistory] = req.History
	}

	if p.avs != nil {
		var memoryVariablesList []string
		for k, v := range p.avs {
			variables[k] = v
			memoryVariablesList = append(memoryVariablesList, fmt.Sprintf("%s: %s\n", k, v))
		}
		variables[placeholderOfVariables] = memoryVariablesList
	}

	return variables, nil
}
