package application

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/model/modelapi"
	"github.com/kiosk404/airi-go/backend/modules/llm/application/convert"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	modelmgrservice "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/dto"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type ModelManagerApplicationService struct {
	appContext *ServiceComponents
	DomainSVC  modelmgrservice.ModelManager
}

func newApplicationService(c *ServiceComponents, domain modelmgrservice.ModelManager) *ModelManagerApplicationService {
	return &ModelManagerApplicationService{
		appContext: c,
		DomainSVC:  domain,
	}
}

func (s *ModelManagerApplicationService) CreateModel(ctx context.Context, req modelapi.CreateModelReq) (int64, error) {
	modelClass := dto.ModelClassDto(req.GetModelClass())
	connection := dto.ModelConnectionDto(req.GetConnection())
	modelShowName := req.GetModelName()
	extra := dto.ModelExtraDto(req.GetEnableBase64URL())
	modelBuilder, err := NewModelBuilder(modelClass, &modelmgr.Model{
		EnableBase64URL: req.EnableBase64URL,
		Connection:      connection,
	})
	if err != nil {
		return 0, err
	}

	logs.DebugX(pkg.ModelName, "create model req: %s, conn: %s", conv.DebugJsonToStr(req), conv.DebugJsonToStr(req.Connection.BaseConnInfo))
	chatModel, err := modelBuilder.Build(ctx, &modelmgr.LLMParams{EnableThinking: ptr.Of(false)})
	if err != nil {
		return 0, err
	}

	respMsgs, err := chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("1+1=?,Just answer with a number, no explanation.")})
	if err != nil {
		return 0, err
	}

	logs.DebugX(pkg.ModelName, "chatModel.Generate resp : %s", conv.DebugJsonToStr(respMsgs))

	id, err := s.DomainSVC.CreateLLMModel(ctx, convert.ModelClassDao(modelClass), modelShowName, convert.ConnectionDao(connection), extra)
	if err != nil {
		return 0, err
	}
	return id, err
}
