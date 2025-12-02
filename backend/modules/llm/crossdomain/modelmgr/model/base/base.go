package base

import (
	"context"
	"errors"
	"os"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model/config"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model/consts"
	"github.com/kiosk404/airi-go/backend/pkg/kvstore"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ternary"
	"gorm.io/gorm"
)

const (
	baseConfigKey = "basic_config"
)

type BaseConfig struct {
	base *kvstore.KVStore[config.BasicConfiguration]
}

func NewBaseConfig(db *gorm.DB) *BaseConfig {
	return &BaseConfig{
		base: kvstore.New[config.BasicConfiguration](db),
	}
}

func (c *BaseConfig) GetBaseConfig(ctx context.Context) (*config.BasicConfiguration, error) {
	conf, err := c.base.Get(ctx, consts.BaseConfigNameSpace, baseConfigKey)
	if err != nil {
		if errors.Is(err, kvstore.ErrKeyNotFound) {
			return getBasicConfigurationFromOldConfig(), nil
		}
	}

	return conf, nil
}

func (c *BaseConfig) SaveBaseConfig(ctx context.Context, v *config.BasicConfiguration) error {
	return c.base.Save(ctx, consts.BaseConfigNameSpace, baseConfigKey, v)
}

func getBasicConfigurationFromOldConfig() *config.BasicConfiguration {
	runnerTypeStr := os.Getenv(consts.CodeRunnerType)
	codeRunnerType := ternary.IFElse(runnerTypeStr == "sandbox", config.CoderunnerTypeSandbox, config.CoderunnerTypeLocal)
	timeoutSecondsStr := os.Getenv(consts.CodeRunnerTimeoutSeconds)
	timeoutSeconds := conv.StrToFloat64D(timeoutSecondsStr, 60)
	memoryLimitMbStr := os.Getenv(consts.CodeRunnerMemoryLimitMB)
	memoryLimitMB := conv.StrToInt64D(memoryLimitMbStr, 100)

	return &config.BasicConfiguration{
		CodeRunnerType: codeRunnerType,
		SandboxConfig: &config.SandboxConfig{
			AllowEnv:       os.Getenv(consts.CodeRunnerAllowEnv),
			AllowRead:      os.Getenv(consts.CodeRunnerAllowRead),
			AllowWrite:     os.Getenv(consts.CodeRunnerAllowWrite),
			AllowNet:       os.Getenv(consts.CodeRunnerAllowNet),
			AllowRun:       os.Getenv(consts.CodeRunnerAllowRun),
			AllowFfi:       os.Getenv(consts.CodeRunnerAllowFFI),
			NodeModulesDir: os.Getenv(consts.CodeRunnerNodeModulesDir),
			TimeoutSeconds: timeoutSeconds,
			MemoryLimitMb:  memoryLimitMB,
		},
	}
}
