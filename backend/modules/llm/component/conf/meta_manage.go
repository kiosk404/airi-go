package conf

import (
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type IConfigMetaManage interface {
	// GetModelMetaConf 获取模型元数据配置
	GetModelMeta(modelClass model.ModelClass, modelName string) (*model.ModelMeta, error)
}
