package convertor

import (
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/common"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func ScenarioDO2DTO(s entity.Scenario) common.Scenario {
	return common.Scenario(s)
}

func ScenarioPtrDTO2DTO(s *common.Scenario) *entity.Scenario {
	if s == nil {
		return nil
	}
	return ptr.Of(entity.Scenario(*s))
}
