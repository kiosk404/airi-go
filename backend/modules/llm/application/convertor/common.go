package convertor

import (
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/common"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func ScenarioDO2DTO(s model.Scenario) common.Scenario {
	return common.Scenario(s)
}

func ScenarioDTO2DO(s common.Scenario) model.Scenario {
	return model.Scenario(s)
}

func ScenarioPtrDTO2DTO(s *common.Scenario) *model.Scenario {
	if s == nil {
		return nil
	}
	return ptr.Of(model.Scenario(*s))
}

func ScenarioPtrDTO2DO(s *common.Scenario) *model.Scenario {
	if s == nil {
		return nil
	}
	return ptr.Of(model.Scenario(*s))
}
