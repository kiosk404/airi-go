package chatmodel

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	acl_openai "github.com/cloudwego/eino-ext/libs/acl/openai"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/ollama/ollama/api"
	ollamaapi "github.com/ollama/ollama/api"
	"google.golang.org/genai"
)

type Builder func(ctx context.Context, config *Config) (ToolCallingChatModel, error)

func NewDefaultFactory() Factory {
	return NewFactory(nil)
}

func NewFactory(customFactory map[entity.Protocol]Builder) Factory {
	protocol2Builder := map[entity.Protocol]Builder{
		entity.ProtocolOpenAI:   openAIBuilder,
		entity.ProtocolClaude:   claudeBuilder,
		entity.ProtocolDeepseek: deepseekBuilder,
		entity.ProtocolGemini:   geminiBuilder,
		entity.ProtocolOllama:   ollamaBuilder,
		entity.ProtocolQwen:     qwenBuilder,
	}

	for p := range customFactory {
		protocol2Builder[p] = customFactory[p]
	}

	return &defaultFactory{protocol2Builder: protocol2Builder}
}

type defaultFactory struct {
	protocol2Builder map[entity.Protocol]Builder
}

func (f *defaultFactory) CreateChatModel(ctx context.Context, protocol entity.Protocol, config *Config) (ToolCallingChatModel, error) {
	if config == nil {
		return nil, fmt.Errorf("[CreateChatModel] config not provided")
	}

	builder, found := f.protocol2Builder[protocol]
	if !found {
		return nil, fmt.Errorf("[CreateChatModel] protocol not support, protocol=%s", protocol)
	}

	return builder(ctx, config)
}

func (f *defaultFactory) SupportProtocol(protocol entity.Protocol) bool {
	_, found := f.protocol2Builder[protocol]
	return found
}

func openAIBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &openai.ChatModelConfig{
		APIKey:           config.GetAPIKey(),
		BaseURL:          config.GetBaseURL(),
		Model:            config.GetModel(),
		MaxTokens:        config.MaxTokens,
		Temperature:      config.Temperature,
		TopP:             config.TopP,
		Stop:             config.Stop,
		PresencePenalty:  config.PresencePenalty,
		FrequencyPenalty: config.FrequencyPenalty,
	}
	if config.Timeout != 0 {
		cfg.Timeout = time.Duration(config.Timeout) * time.Millisecond
	}

	configProtocol := config.ProtocolConfigOpenai

	if configProtocol != nil {
		cfg.ByAzure = configProtocol.GetByAzure()
		cfg.APIVersion = configProtocol.GetAPIVersion()
		var js acl_openai.ChatCompletionResponseFormatJSONSchema
		if configProtocol.GetResponseFormatType() != "" {
			if err := sonic.UnmarshalString(configProtocol.GetResponseFormatJSONSchema(), js); err != nil {
				return nil, err
			}
		}
		cfg.ResponseFormat = &acl_openai.ChatCompletionResponseFormat{
			Type:       acl_openai.ChatCompletionResponseFormatType(configProtocol.GetResponseFormatType()),
			JSONSchema: &js,
		}

	}
	return openai.NewChatModel(ctx, cfg)
}

func claudeBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &claude.Config{
		APIKey:        config.GetAPIKey(),
		Model:         config.GetModel(),
		Temperature:   config.Temperature,
		TopP:          config.TopP,
		StopSequences: config.Stop,
	}
	if ptr.From(config.BaseURL) != "" {
		cfg.BaseURL = config.BaseURL
	}
	if config.MaxTokens != nil {
		cfg.MaxTokens = *config.MaxTokens
	}
	if config.TopK != nil {
		cfg.TopK = ptr.Of(int32(*config.TopK))
	}
	configProtocol := config.ProtocolConfigClaude
	if configProtocol != nil {
		cfg.ByBedrock = configProtocol.GetByBedrock()
		cfg.AccessKey = configProtocol.GetAccessKey()
		cfg.SecretAccessKey = configProtocol.GetSecretAccessKey()
		cfg.SessionToken = configProtocol.GetSessionToken()
		cfg.Region = configProtocol.GetRegion()
	}
	if config.EnableThinking != nil {
		cfg.Thinking = &claude.Thinking{
			Enable: ptr.From(config.EnableThinking),
		}
		if configProtocol != nil && configProtocol.BudgetTokens != nil {
			cfg.Thinking.BudgetTokens = int(configProtocol.GetBudgetTokens())
		}
	}
	return claude.NewChatModel(ctx, cfg)
}

func deepseekBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &deepseek.ChatModelConfig{
		APIKey:  config.GetAPIKey(),
		BaseURL: config.GetBaseURL(),
		Model:   config.GetModel(),
		Stop:    config.Stop,
	}
	if config.Timeout != 0 {
		cfg.Timeout = time.Duration(config.Timeout) * time.Millisecond
	}
	if config.Temperature != nil {
		cfg.Temperature = *config.Temperature
	}
	if config.FrequencyPenalty != nil {
		cfg.FrequencyPenalty = *config.FrequencyPenalty
	}
	if config.PresencePenalty != nil {
		cfg.PresencePenalty = *config.PresencePenalty
	}
	if config.MaxTokens != nil {
		cfg.MaxTokens = *config.MaxTokens
	}
	if config.TopP != nil {
		cfg.TopP = *config.TopP
	}

	protocolConfig := config.ProtocolConfigDeepseek
	if protocolConfig != nil {
		cfg.ResponseFormatType = deepseek.ResponseFormatType(protocolConfig.GetResponseFormatType())
	}

	return deepseek.NewChatModel(ctx, cfg)
}

func ollamaBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	configProtocol := config.ProtocolConfigOllama

	cfg := &ollama.ChatModelConfig{
		BaseURL:    config.GetBaseURL(),
		Model:      config.GetModel(),
		HTTPClient: nil,
		Format:     nil,
		KeepAlive:  nil,
		Options: &api.Options{
			TopK:             ptr.From(config.TopK),
			TopP:             ptr.From(config.TopP),
			Temperature:      ptr.From(config.Temperature),
			PresencePenalty:  ptr.From(config.PresencePenalty),
			FrequencyPenalty: ptr.From(config.FrequencyPenalty),
			Stop:             config.Stop,
		},
	}
	if config.Timeout != 0 {
		cfg.Timeout = time.Duration(config.Timeout) * time.Millisecond
	}
	if config.EnableThinking != nil {
		cfg.Thinking = &ollamaapi.ThinkValue{Value: ptr.From(config.EnableThinking)}
	}
	if configProtocol.IsSetKeepAliveMs() {
		cfg.KeepAlive = ptr.Of(time.Duration(configProtocol.GetKeepAliveMs()) * time.Millisecond)
	}
	if configProtocol.IsSetFormat() {
		cfg.Format = []byte(configProtocol.GetFormat())
	}
	return ollama.NewChatModel(ctx, cfg)
}

func qwenBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &qwen.ChatModelConfig{
		APIKey:           config.GetAPIKey(),
		BaseURL:          config.GetBaseURL(),
		Model:            config.GetModel(),
		MaxTokens:        config.MaxTokens,
		Temperature:      config.Temperature,
		TopP:             config.TopP,
		Stop:             config.Stop,
		PresencePenalty:  config.PresencePenalty,
		FrequencyPenalty: config.FrequencyPenalty,
		EnableThinking:   config.EnableThinking,
	}
	if config.Timeout != 0 {
		cfg.Timeout = time.Duration(config.Timeout) * time.Millisecond
	}
	configProtocol := config.ProtocolConfigQwen
	if configProtocol != nil {
		if configProtocol.ResponseFormatType != nil {
			var js acl_openai.ChatCompletionResponseFormatJSONSchema
			if configProtocol.ResponseFormatJSONSchema != nil {
				if err := sonic.UnmarshalString(configProtocol.GetResponseFormatJSONSchema(), js); err != nil {
					return nil, err
				}
			}
			cfg.ResponseFormat = &acl_openai.ChatCompletionResponseFormat{
				Type:       acl_openai.ChatCompletionResponseFormatType(ptr.From(configProtocol.ResponseFormatType)),
				JSONSchema: &js,
			}
		}
	}
	return qwen.NewChatModel(ctx, cfg)
}

func geminiBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	gc := &genai.ClientConfig{
		APIKey: config.GetAPIKey(),
		HTTPOptions: genai.HTTPOptions{
			BaseURL: config.GetBaseURL(),
		},
	}
	configProtocol := config.ProtocolConfigGemini

	if configProtocol != nil {
		gc.Backend = genai.Backend(configProtocol.GetBackend())
		gc.Project = configProtocol.GetProject()
		gc.Location = configProtocol.GetLocation()
		gc.HTTPOptions.APIVersion = configProtocol.GetAPIVersion()
		gc.HTTPOptions.Headers = configProtocol.Headers
	}

	client, err := genai.NewClient(ctx, gc)
	if err != nil {
		return nil, err
	}

	cfg := &gemini.Config{
		Client:      client,
		Model:       config.GetModel(),
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
		TopP:        config.TopP,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
			ThinkingBudget:  nil,
		},
	}
	if config.TopK != nil {
		cfg.TopK = ptr.Of(int32(ptr.From(config.TopK)))
	}

	protocolConfig := config.ProtocolConfigGemini
	if protocolConfig != nil {
		if protocolConfig.IncludeThoughts != nil {
			cfg.ThinkingConfig.IncludeThoughts = ptr.From(protocolConfig.IncludeThoughts)
		}
		if protocolConfig.ThinkingBudget != nil {
			cfg.ThinkingConfig.ThinkingBudget = protocolConfig.ThinkingBudget
		}
	}

	cm, err := gemini.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return cm, nil
}
