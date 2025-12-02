package service

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	confpkg "github.com/kiosk404/airi-go/backend/pkg/conf"
	"github.com/kiosk404/airi-go/backend/pkg/conf/viper"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type ModelMetaConf struct {
	Provider2Models map[string]map[string]ModelMeta `thrift:"provider2models,2" form:"provider2models" json:"provider2models" query:"provider2models"`
}

type ModelMeta model.ModelMeta

var modelMetaConf *ModelMetaConf

func ModelMetaConfFactory(appRootPath string) (confpkg.IConfigLoaderFactory, error) {
	var configFactory confpkg.IConfigLoaderFactory
	var err error
	projectConfigDir := fmt.Sprintf("%s/conf", appRootPath)
	configOptionList := []viper.FileConfLoaderFactoryOpt{viper.WithFactoryConfigPath(projectConfigDir)}
	if configFactory, err = viper.NewFileConfigLoaderFactory(configOptionList...); err != nil {
		return nil, fmt.Errorf("init config loader factory, err=%w", err)
	}

	return configFactory, nil
}

// 当前的模型元数据配置，从model_meta.json文件中加载, 后面从数据库中加载
func initModelMetaConf(ctx context.Context, factory confpkg.IConfigLoaderFactory) (*ModelMetaConf, error) {
	loader, err := factory.NewConfigLoader("model_meta.json")
	if err != nil {
		return nil, fmt.Errorf("error reading model_meta.json: %w", err)
	}

	err = loader.Unmarshal(ctx, &modelMetaConf)

	if err != nil {
		return nil, fmt.Errorf("error Unmarshal model_meta.json: %w", err)
	}

	return modelMetaConf, nil
}

func (c *ModelMetaConf) GetModelMeta(modelClass model.ModelClass, modelName string) (*model.ModelMeta, error) {
	modelName2Meta, ok := c.Provider2Models[modelClass.String()]
	if !ok {
		return nil, fmt.Errorf("model meta not found for model class %v", modelClass)
	}

	modelMeta, ok := modelName2Meta[modelName]
	if ok {
		logs.InfoX(pkg.ModelName, "get model meta for model class %v and model name %v", modelClass, modelName)
		return deepCopyModelMeta(&modelMeta)
	}

	const defaultKey = "default"
	modelMeta, ok = modelName2Meta[defaultKey]
	if ok {
		logs.InfoX(pkg.ModelName, "use default model meta for model class %v and model name %v", modelClass, modelName)
		return deepCopyModelMeta(&modelMeta)
	}

	return nil, fmt.Errorf("model meta not found for model class %v and model name %v", modelClass, modelName)
}

func deepCopyModelMeta(meta *ModelMeta) (*model.ModelMeta, error) {
	if meta == nil {
		return nil, nil
	}
	newObj := &model.ModelMeta{}
	err := copier.CopyWithOption(newObj, meta, copier.Option{DeepCopy: true, IgnoreEmpty: true})
	if err != nil {
		return nil, fmt.Errorf("error copy model meta: %w", err)
	}

	return newObj, nil
}
