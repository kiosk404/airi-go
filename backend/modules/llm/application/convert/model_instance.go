package convert

import (
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func ModelInstance(model *entity.ModelInstance) *modelmgr.Model {
	return &modelmgr.Model{
		ID:              model.ID,
		Provider:        ptr.Of(model.Provider),
		DisplayInfo:     ptr.Of(model.DisplayInfo),
		Capability:      ptr.Of(model.Capability),
		Connection:      ptr.Of(model.Connection.Model()),
		Type:            model.Type.Model(),
		Parameters:      model.Parameters,
		EnableBase64URL: model.Extra.EnableBase64URL,
	}
}

func ModelClassDao(modelClass modelmgr.ModelClass) entity.ModelClass {
	return entity.ModelClass{ModelClass: ptr.Of(modelClass)}
}

func ConnectionDao(connection *modelmgr.Connection) *entity.Connection {
	return &entity.Connection{Connection: connection}
}
