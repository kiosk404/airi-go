package application

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"google.golang.org/genai"
)

type geminiModelBuilder struct {
	cfg *model.Model
}

func newGeminiModelBuilder(cfg *model.Model) Service {
	return &geminiModelBuilder{
		cfg: cfg,
	}
}

func (g *geminiModelBuilder) getDefaultGeminiConfig() *gemini.Config {
	return &gemini.Config{}
}

func (g *geminiModelBuilder) getDefaultGenaiConfig() *genai.ClientConfig {
	return &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{
			BaseURL: "https://generativelanguage.googleapis.com/",
		},
	}
}

func (g *geminiModelBuilder) applyParamsToGeminiConfig(conf *gemini.Config, params *model.LLMParams) {
	if params == nil {
		return
	}

	conf.TopK = params.TopK
	conf.TopP = params.TopP

	if params.Temperature != nil {
		conf.Temperature = ptr.Of(*params.Temperature)
	}

	if params.MaxTokens != 0 {
		conf.MaxTokens = ptr.Of(params.MaxTokens)
	}

	if params.EnableThinking != nil {
		conf.ThinkingConfig = &genai.ThinkingConfig{
			IncludeThoughts: *params.EnableThinking,
		}
	}
}

func (g *geminiModelBuilder) Build(ctx context.Context, params *model.LLMParams) (ToolCallingChatModel, error) {
	base := g.cfg.Connection.BaseConnInfo

	clientCfg := g.getDefaultGenaiConfig()
	if base.BaseURL != "" {
		clientCfg.HTTPOptions.BaseURL = base.BaseURL
	}

	clientCfg.APIKey = base.APIKey
	if g.cfg.Connection.Gemini != nil {
		clientCfg.Backend = genai.Backend(g.cfg.Connection.Gemini.Backend)
		clientCfg.Project = g.cfg.Connection.Gemini.Project
		clientCfg.Location = g.cfg.Connection.Gemini.Location
	}

	client, err := genai.NewClient(ctx, clientCfg)
	if err != nil {
		return nil, err
	}

	conf := g.getDefaultGeminiConfig()
	conf.Client = client
	conf.Model = base.Model

	switch base.ThinkingType {
	case model.ThinkingType_Enable:
		conf.ThinkingConfig = &genai.ThinkingConfig{
			IncludeThoughts: true,
		}
	case model.ThinkingType_Disable:
		conf.ThinkingConfig = &genai.ThinkingConfig{
			IncludeThoughts: false,
		}
	}

	g.applyParamsToGeminiConfig(conf, params)

	return gemini.NewChatModel(ctx, conf)
}
