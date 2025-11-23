package model

type LLMParams struct {
	ModelName         string             `json:"modelName"`
	ModelType         int64              `json:"modelType"`
	Prompt            string             `json:"prompt"` // user prompt
	Temperature       *float64           `json:"temperature"`
	FrequencyPenalty  float64            `json:"frequencyPenalty"`
	PresencePenalty   float64            `json:"presencePenalty"`
	MaxTokens         int                `json:"maxTokens"`
	TopP              *float64           `json:"topP"`
	TopK              *int               `json:"topK"`
	EnableChatHistory bool               `json:"enableChatHistory"`
	SystemPrompt      string             `json:"systemPrompt"`
	ResponseFormat    ResponseFormatType `json:"responseFormat"`
	ChatHistoryRound  int64              `json:"chatHistoryRound"`
}

type ResponseFormatType int64

const (
	ResponseFormatText     ResponseFormatType = 0
	ResponseFormatMarkdown ResponseFormatType = 1
	ResponseFormatJSON     ResponseFormatType = 2
)
