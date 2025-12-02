package config

import (
	"fmt"
)

type KnowledgeConfig struct {
	EmbeddingConfig *EmbeddingConfig `json:"embedding_config" query:"embedding_config"`
	RerankConfig    *RerankConfig    `json:"rerank_config" query:"rerank_config"`
	OcrConfig       *OCRConfig       `json:"ocr_config" query:"ocr_config"`
	ParserConfig    *ParserConfig    `json:"parser_config" query:"parser_config"`
	BuiltinModelID  int64            `json:"builtin_model_id" query:"builtin_model_id"`
}

type EmbeddingConfig struct {
	Type         EmbeddingType        `json:"type" query:"type"`
	MaxBatchSize int32                `json:"max_batch_size" query:"max_batch_size"`
	Connection   *EmbeddingConnection `json:"connection" query:"connection"`
}

type RerankConfig struct {
	Type           RerankType      `json:"type" query:"type"`
	VikingDBConfig *VikingDBConfig `json:"vikingdb_config" query:"vikingdb_config"`
}

func NewRerankConfig() *RerankConfig {
	return &RerankConfig{}
}

func (p *RerankConfig) InitDefault() {
}

func (p *RerankConfig) GetType() (v RerankType) {
	return p.Type
}

type VikingDBConfig struct {
	Ak     string `json:"ak" query:"ak"`
	Sk     string `json:"sk" query:"sk"`
	Host   string `json:"host" query:"host"`
	Region string `json:"region" query:"region"`
	Model  string `json:"model" query:"model"`
}

type OCRConfig struct {
	Type            OCRType `json:"type" query:"type"`
	VolcengineAk    string  `json:"volcengine_ak" query:"volcengine_ak"`
	VolcengineSk    string  `json:"volcengine_sk" query:"volcengine_sk"`
	PaddleocrAPIURL string  `json:"paddleocr_api_url" query:"paddleocr_api_url"`
}

func NewOCRConfig() *OCRConfig {
	return &OCRConfig{}
}

func (p *OCRConfig) InitDefault() {
}

func (p *OCRConfig) GetType() (v OCRType) {
	return p.Type
}

func (p *OCRConfig) GetVolcengineAk() (v string) {
	return p.VolcengineAk
}

func (p *OCRConfig) GetVolcengineSk() (v string) {
	return p.VolcengineSk
}

func (p *OCRConfig) GetPaddleocrAPIURL() (v string) {
	return p.PaddleocrAPIURL
}

type ParserConfig struct {
	Type                     ParserType `json:"type" query:"type"`
	PaddleocrStructureAPIURL string     `json:"paddleocr_structure_api_url" query:"paddleocr_structure_api_url"`
}

func NewParserConfig() *ParserConfig {
	return &ParserConfig{}
}

func (p *ParserConfig) InitDefault() {
}

func (p *ParserConfig) GetType() (v ParserType) {
	return p.Type
}

func (p *ParserConfig) GetPaddleocrStructureAPIURL() (v string) {
	return p.PaddleocrStructureAPIURL
}

type OCRType int64

const (
	OCRType_Volcengine OCRType = 0
	OCRType_Paddleocr  OCRType = 1
)

func (p OCRType) String() string {
	switch p {
	case OCRType_Volcengine:
		return "Volcengine"
	case OCRType_Paddleocr:
		return "Paddleocr"
	}
	return "<UNSET>"
}

func OCRTypeFromString(s string) (OCRType, error) {
	switch s {
	case "Volcengine":
		return OCRType_Volcengine, nil
	case "Paddleocr":
		return OCRType_Paddleocr, nil
	}
	return OCRType(0), fmt.Errorf("not a valid OCRType string")
}

type RerankType int64

const (
	RerankType_VikingDB RerankType = 0
	RerankType_RRF      RerankType = 1
)

func (p RerankType) String() string {
	switch p {
	case RerankType_VikingDB:
		return "VikingDB"
	case RerankType_RRF:
		return "RRF"
	}
	return "<UNSET>"
}

type EmbeddingType int64

const (
	EmbeddingType_Ark    EmbeddingType = 0
	EmbeddingType_OpenAI EmbeddingType = 1
	EmbeddingType_Ollama EmbeddingType = 2
	EmbeddingType_Gemini EmbeddingType = 3
	EmbeddingType_HTTP   EmbeddingType = 4
)

func (p EmbeddingType) String() string {
	switch p {
	case EmbeddingType_Ark:
		return "Ark"
	case EmbeddingType_OpenAI:
		return "OpenAI"
	case EmbeddingType_Ollama:
		return "Ollama"
	case EmbeddingType_Gemini:
		return "Gemini"
	case EmbeddingType_HTTP:
		return "HTTP"
	}
	return "<UNSET>"
}

func EmbeddingTypeFromString(s string) (EmbeddingType, error) {
	switch s {
	case "Ark":
		return EmbeddingType_Ark, nil
	case "OpenAI":
		return EmbeddingType_OpenAI, nil
	case "Ollama":
		return EmbeddingType_Ollama, nil
	case "Gemini":
		return EmbeddingType_Gemini, nil
	case "HTTP":
		return EmbeddingType_HTTP, nil
	}
	return EmbeddingType(0), fmt.Errorf("not a valid EmbeddingType string")
}

type ParserType int64

const (
	ParserType_builtin   ParserType = 0
	ParserType_Paddleocr ParserType = 1
)

func (p ParserType) String() string {
	switch p {
	case ParserType_builtin:
		return "builtin"
	case ParserType_Paddleocr:
		return "Paddleocr"
	}
	return "<UNSET>"
}

func ParserTypeFromString(s string) (ParserType, error) {
	switch s {
	case "builtin":
		return ParserType_builtin, nil
	case "Paddleocr":
		return ParserType_Paddleocr, nil
	}
	return ParserType(0), fmt.Errorf("not a valid ParserType string")
}

func (p *KnowledgeConfig) GetEmbeddingConfig() (v *EmbeddingConfig) {
	if p.EmbeddingConfig == nil {
		return nil
	}

	return p.EmbeddingConfig
}
