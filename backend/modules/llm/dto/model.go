package dto

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/api/model/modelapi"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
)

func ModelClassDto(class developer_api.ModelClass) modelmgr.ModelClass {
	return modelmgr.ModelClass(class)
}

func ModelConnectionDto(connection *modelapi.Connection) *modelmgr.Connection {
	return &modelmgr.Connection{
		BaseConnInfo: &modelmgr.BaseConnectionInfo{
			BaseURL:      connection.BaseConnInfo.BaseURL,
			APIKey:       connection.BaseConnInfo.APIKey,
			Model:        connection.BaseConnInfo.Model,
			ThinkingType: modelmgr.ThinkingType(connection.BaseConnInfo.ThinkingType),
		},
	}
}

func ModelExtraDto(extraEnableBase64URL bool) *entity.ModelExtra {
	return &entity.ModelExtra{
		ModelExtra: &modelmgr.ModelExtra{
			EnableBase64URL: extraEnableBase64URL,
		},
	}
}
