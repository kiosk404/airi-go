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
		APIKey:           config.APIKey,
		BaseURL:          config.BaseURL,
		Model:            config.Model,
		MaxTokens:        config.MaxTokens,
		Temperature:      config.Temperature,
		TopP:             config.TopP,
		Stop:             config.Stop,
		PresencePenalty:  config.PresencePenalty,
		FrequencyPenalty: config.FrequencyPenalty,
	}
	if config.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*config.TimeoutMs) * time.Millisecond
	}

	configProtocol := config.ProtocolConfigOpenAI

	if configProtocol != nil {
		cfg.ByAzure = configProtocol.ByAzure
		cfg.APIVersion = configProtocol.ApiVersion
		var js acl_openai.ChatCompletionResponseFormatJSONSchema
		if configProtocol.ResponseFormatJsonSchema != "" {
			if err := sonic.UnmarshalString(configProtocol.ResponseFormatJsonSchema, js); err != nil {
				return nil, err
			}
		}
		cfg.ResponseFormat = &acl_openai.ChatCompletionResponseFormat{
			Type:       acl_openai.ChatCompletionResponseFormatType(configProtocol.ResponseFormatType),
			JSONSchema: &js,
		}

	}
	return openai.NewChatModel(ctx, cfg)
}

func claudeBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &claude.Config{
		APIKey:        config.APIKey,
		Model:         config.Model,
		Temperature:   config.Temperature,
		TopP:          config.TopP,
		StopSequences: config.Stop,
	}
	if config.BaseURL != "" {
		cfg.BaseURL = &config.BaseURL
	}
	if config.MaxTokens != nil {
		cfg.MaxTokens = *config.MaxTokens
	}
	if config.TopK != nil {
		cfg.TopK = ptr.Of(int32(*config.TopK))
	}
	configProtocol := config.ProtocolConfigClaude
	if configProtocol != nil {
		cfg.ByBedrock = configProtocol.ByBedrock
		cfg.AccessKey = configProtocol.AccessKey
		cfg.SecretAccessKey = configProtocol.SecretAccessKey
		cfg.SessionToken = configProtocol.SessionToken
		cfg.Region = configProtocol.Region
	}
	if config.EnableThinking != nil {
		cfg.Thinking = &claude.Thinking{
			Enable: ptr.From(config.EnableThinking),
		}
		if configProtocol != nil && configProtocol.BudgetTokens != nil {
			cfg.Thinking.BudgetTokens = ptr.From(configProtocol.BudgetTokens)
		}
	}
	return claude.NewChatModel(ctx, cfg)
}

func deepseekBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &deepseek.ChatModelConfig{
		APIKey:  config.APIKey,
		BaseURL: config.BaseURL,
		Model:   config.Model,
		Stop:    config.Stop,
	}
	if config.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*config.TimeoutMs) * time.Millisecond
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

	protocolConfig := config.ProtocolConfigDeepSeek
	if protocolConfig != nil {
		cfg.ResponseFormatType = deepseek.ResponseFormatType(protocolConfig.ResponseFormatType)
	}

	return deepseek.NewChatModel(ctx, cfg)
}

func ollamaBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &ollama.ChatModelConfig{
		BaseURL:    config.BaseURL,
		HTTPClient: nil,
		Model:      config.Model,
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
	if config.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*config.TimeoutMs) * time.Millisecond
	}
	if config.EnableThinking != nil {
		cfg.Thinking = &ollamaapi.ThinkValue{Value: ptr.From(config.EnableThinking)}
	}
	return ollama.NewChatModel(ctx, cfg)
}

func qwenBuilder(ctx context.Context, config *Config) (ToolCallingChatModel, error) {
	cfg := &qwen.ChatModelConfig{
		APIKey:           config.APIKey,
		BaseURL:          config.BaseURL,
		Model:            config.Model,
		MaxTokens:        config.MaxTokens,
		Temperature:      config.Temperature,
		TopP:             config.TopP,
		Stop:             config.Stop,
		PresencePenalty:  config.PresencePenalty,
		FrequencyPenalty: config.FrequencyPenalty,
		EnableThinking:   config.EnableThinking,
	}
	if config.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*config.TimeoutMs) * time.Millisecond
	}
	configProtocol := config.ProtocolConfigQwen
	if configProtocol != nil {
		if configProtocol.ResponseFormatType != nil {
			var js acl_openai.ChatCompletionResponseFormatJSONSchema
			if configProtocol.ResponseFormatJsonSchema != nil {
				if err := sonic.UnmarshalString(ptr.From(configProtocol.ResponseFormatJsonSchema), js); err != nil {
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
		APIKey: config.APIKey,
		HTTPOptions: genai.HTTPOptions{
			BaseURL: config.BaseURL,
		},
	}
	configProtocol := config.ProtocolConfigGemini

	if configProtocol != nil {
		gc.Backend = configProtocol.Backend
		gc.Project = configProtocol.Project
		gc.Location = configProtocol.Location
		gc.HTTPOptions.APIVersion = configProtocol.APIVersion
		gc.HTTPOptions.Headers = configProtocol.Headers

	}

	client, err := genai.NewClient(ctx, gc)
	if err != nil {
		return nil, err
	}

	cfg := &gemini.Config{
		Client:      client,
		Model:       config.Model,
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

		if protocolConfig.ResponseSchema != nil && *protocolConfig.ResponseSchema != "" {
			if err := sonic.UnmarshalString(*protocolConfig.ResponseSchema, &cfg.ResponseSchema); err != nil {
				return nil, err
			}
		}
		cfg.EnableCodeExecution = protocolConfig.EnableCodeExecution
		for _, ss := range protocolConfig.SafetySettings {
			cfg.SafetySettings = append(cfg.SafetySettings, &genai.SafetySetting{
				Category:  genai.HarmCategory(ss.Category),
				Threshold: genai.HarmBlockThreshold(ss.Threshold),
			})
		}
	}

	cm, err := gemini.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return cm, nil
}
