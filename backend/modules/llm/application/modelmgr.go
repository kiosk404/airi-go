package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
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

func (s *ModelManagerApplicationService) GetProviderModelList(ctx context.Context, modelType modelapi.ModelType) ([]*modelapi.ProviderModelList, error) {
	modelProviderList := getModelProviderList()
	resp := make([]*modelapi.ProviderModelList, 0, len(modelProviderList))

	allModels, err := s.DomainSVC.ListModelByType(ctx, dto.ModelTypeDto(modelType), 30)
	if err != nil {
		return nil, err
	}

	modelClass2Models := make(map[developer_api.ModelClass][]*modelapi.Model)
	for _, model := range allModels {
		m := convert.ToModel(ctx, s.appContext.TosClient, model)
		modelClass2Models[m.Provider.ModelClass] = append(modelClass2Models[m.Provider.ModelClass], m)
		if m.Connection != nil && m.Connection.BaseConnInfo != nil {
			apiKey := m.Connection.BaseConnInfo.APIKey
			if apiKey != "" {
				n := len(apiKey)
				if n <= 4 {
					m.Connection.BaseConnInfo.APIKey = strings.Repeat("*", n)
				} else if n <= 8 {
					m.Connection.BaseConnInfo.APIKey = fmt.Sprintf("%s***%s", apiKey[:2], apiKey[n-2:])
				} else {
					m.Connection.BaseConnInfo.APIKey = fmt.Sprintf("%s***%s", apiKey[:4], apiKey[n-4:])
				}
			}
		}
	}

	for _, provider := range modelProviderList {
		if provider.IconURI != "" {
			url, err := s.appContext.TosClient.GetObjectUrl(ctx, provider.IconURI)
			if err != nil {
				logs.WarnX(pkg.ModelName, "get model icon url failed, err: %v", err)
			} else {
				provider.IconURL = url
			}
		}
		resp = append(resp, &modelapi.ProviderModelList{
			Provider:  provider,
			ModelList: modelClass2Models[provider.ModelClass],
		})
	}

	return resp, nil
}

func (s *ModelManagerApplicationService) GetInUseModelList(ctx context.Context, modelType modelapi.ModelType) ([]*modelapi.ProviderModelList, error) {
	resp := make([]*modelapi.ProviderModelList, 0, 0)

	allModels, err := s.DomainSVC.ListModelByType(ctx, dto.ModelTypeDto(modelType), 30)
	if err != nil {
		return nil, err
	}

	modelClass2Models := make(map[developer_api.ModelClass][]*modelapi.Model)
	for _, model := range allModels {
		m := convert.ToModel(ctx, s.appContext.TosClient, model)
		modelClass2Models[m.Provider.ModelClass] = append(modelClass2Models[m.Provider.ModelClass], m)
		if m.Connection != nil && m.Connection.BaseConnInfo != nil {
			apiKey := m.Connection.BaseConnInfo.APIKey
			if apiKey != "" {
				n := len(apiKey)
				if n <= 4 {
					m.Connection.BaseConnInfo.APIKey = strings.Repeat("*", n)
				} else if n <= 8 {
					m.Connection.BaseConnInfo.APIKey = fmt.Sprintf("%s***%s", apiKey[:2], apiKey[n-2:])
				} else {
					m.Connection.BaseConnInfo.APIKey = fmt.Sprintf("%s***%s", apiKey[:4], apiKey[n-4:])
				}
			}
		}
		if m.Provider.IconURI != "" {
			url, err := s.appContext.TosClient.GetObjectUrl(ctx, m.Provider.IconURI)
			if err != nil {
				logs.WarnX(pkg.ModelName, "get model icon url failed, err: %v", err)
			} else {
				m.Provider.IconURL = url
			}
		}
		resp = append(resp, &modelapi.ProviderModelList{
			Provider:  m.Provider,
			ModelList: modelClass2Models[m.Provider.ModelClass],
		})
	}

	return resp, nil
}
