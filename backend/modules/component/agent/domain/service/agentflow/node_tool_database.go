package agentflow

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type databaseConfig struct {
	agentIdentity *entity.AgentIdentity
	userID        string
	spaceID       int64

	databaseConf []*bot_common.Database
}

type databaseTool struct {
	agentIdentity *entity.AgentIdentity
	connectorUID  string
	spaceID       int64

	databaseID int64

	name           string
	promptDisabled bool
}

type ExecuteSQLRequest struct {
	SQL string `json:"sql" jsonschema:"description=SQL query to execute against the database. You can use standard SQL syntax like SELECT, INSERT, UPDATE, DELETE."`
}

func (d *databaseTool) Invoke(ctx context.Context, req ExecuteSQLRequest) (string, error) {
	if req.SQL == "" {
		return "", fmt.Errorf("sql is empty")
	}
	if d.promptDisabled {
		return "the tool to be called is not available", nil
	}

	panic("implement me")
}

func newDatabaseTools(ctx context.Context, conf *databaseConfig) ([]tool.InvokableTool, error) {
	if conf == nil || len(conf.databaseConf) == 0 {
		return nil, nil
	}

	dbInfos := conf.databaseConf
	tools := make([]tool.InvokableTool, 0, len(dbInfos))
	for _, dbInfo := range dbInfos {
		tID, err := strconv.ParseInt(dbInfo.GetTableId(), 10, 64)
		if err != nil {
			return nil, err
		}
		d := &databaseTool{
			spaceID:        conf.spaceID,
			connectorUID:   conf.userID,
			agentIdentity:  conf.agentIdentity,
			promptDisabled: dbInfo.GetPromptDisabled(),
			name:           dbInfo.GetTableName(),
			databaseID:     tID,
		}

		dbTool, err := utils.InferTool(dbInfo.GetTableName(), buildDatabaseToolDescription(dbInfo), d.Invoke)
		if err != nil {
			return nil, err
		}

		tools = append(tools, dbTool)
	}

	return tools, nil
}

func buildDatabaseToolDescription(tableInfo *bot_common.Database) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Mysql query tool. Table name is '%s'.", tableInfo.GetTableName()))
	if tableInfo.GetTableDesc() != "" {
		sb.WriteString(fmt.Sprintf(" This table's desc is %s.", tableInfo.GetTableDesc()))
	}
	sb.WriteString("\n\nTable structure:\n")

	for _, field := range tableInfo.FieldList {
		if field.Name == nil || field.Type == nil {
			continue
		}

		fieldType := getFieldTypeString(*field.Type)
		sb.WriteString(fmt.Sprintf("- %s (%s)", *field.Name, fieldType))

		if field.Desc != nil && *field.Desc != "" {
			sb.WriteString(fmt.Sprintf(": %s", *field.Desc))
		}

		if field.MustRequired != nil && *field.MustRequired {
			sb.WriteString(" (required)")
		}

		sb.WriteString("\n")
	}

	sb.WriteString("\nUse SQL to query this table. You can write SQL statements directly to operate.")
	return sb.String()
}

func getFieldTypeString(fieldType bot_common.FieldItemType) string {
	switch fieldType {
	case bot_common.FieldItemType_Text:
		return "text"
	case bot_common.FieldItemType_Number:
		return "number"
	case bot_common.FieldItemType_Date:
		return "date"
	case bot_common.FieldItemType_Float:
		return "float"
	case bot_common.FieldItemType_Boolean:
		return "bool"
	default:
		return "invalid"
	}
}
