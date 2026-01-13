package entity

import (
	"time"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type ModelInstance struct {
	ID          int64
	Type        ModelType
	Provider    model.ModelProvider
	DisplayInfo model.DisplayInfo
	IsSelected  bool
	Connection  Connection
	Capability  model.ModelAbility
	Parameters  []model.ModelParameter
	Extra       ModelExtra
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type ModelClass struct{ *model.ModelClass }

func (m *ModelClass) Model() model.ModelClass {
	return ptr.From(m.ModelClass)
}

type ModelType struct{ *model.ModelType }

func (t *ModelType) Model() model.ModelType {
	return ptr.From(t.ModelType)
}

func LLMModelType() ModelType {
	return ModelType{ModelType: ptr.Of(model.ModelType_LLM)}
}

func EmbeddingModelType() ModelType {
	return ModelType{ModelType: ptr.Of(model.ModelType_TextEmbedding)}
}

func RerankModelType() ModelType {
	return ModelType{ModelType: ptr.Of(model.ModelType_Rerank)}
}

type Connection struct{ *model.Connection }

func (c *Connection) Model() model.Connection {
	return ptr.From(c.Connection)
}

type ModelExtra struct{ *model.ModelExtra }

func (e *ModelExtra) Model() model.ModelExtra {
	return ptr.From(e.ModelExtra)
}
