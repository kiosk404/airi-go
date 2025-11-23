package entity

import (
	"time"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type ModelRequestRecord struct {
	ID                  int64          `json:"id"`
	UserID              string         `json:"user_id"`
	UsageScene          model.Scenario `json:"usage_scene"`
	UsageSceneEntityID  string         `json:"usage_scene_entity_id"`
	Protocol            model.Protocol `json:"protocol"`
	ModelIdentification string         `json:"model_identification"`
	ModelAk             string         `json:"model_ak"`
	ModelID             string         `json:"model_id"`
	ModelName           string         `json:"model_name"`
	InputToken          int64          `json:"input_token"`
	OutputToken         int64          `json:"output_token"`
	LogId               string         `json:"log_id"`
	ErrorCode           string         `json:"error_code"`
	ErrorMsg            *string        `json:"error_msg"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}
