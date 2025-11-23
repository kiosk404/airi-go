package agentflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
)

type retrieverConfig struct {
	knowledgeConfig *bot_common.Knowledge
}

func newKnowledgeRetriever(_ context.Context, conf *retrieverConfig) (*knowledgeRetriever, error) {
	return &knowledgeRetriever{
		knowledgeConfig: conf.knowledgeConfig,
	}, nil
}

type knowledgeRetriever struct {
	knowledgeConfig *bot_common.Knowledge
}

func (r *knowledgeRetriever) Retrieve(ctx context.Context, req *AgentRequest) ([]*schema.Document, error) {
	var docs []*schema.Document
	return docs, nil
}

func (r *knowledgeRetriever) PackRetrieveResultInfo(ctx context.Context, docs []*schema.Document) (string, error) {
	packedRes := strings.Builder{}
	for idx, doc := range docs {
		if doc == nil {
			continue
		}
		packedRes.WriteString(fmt.Sprintf("---\nrecall slice %d: %s\n", idx+1, doc.Content))
	}
	return packedRes.String(), nil
}
