package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	pluginentity "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	agentrunentity "github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	mgrentity "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/llmmgr/model"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

var table2Columns2Model = map[string]map[string]any{
	"single_agent_draft": {
		"variable":                   []*bot_common.Variable{},
		"model_info":                 &bot_common.ModelInfo{},
		"onboarding_info":            &bot_common.OnboardingInfo{},
		"prompt":                     &bot_common.PromptInfo{},
		"plugin":                     []*bot_common.PluginInfo{},
		"knowledge":                  &bot_common.Knowledge{},
		"workflow":                   []*bot_common.WorkflowInfo{},
		"suggest_reply":              &bot_common.SuggestReplyInfo{},
		"jump_config":                &bot_common.JumpConfig{},
		"background_image_info_list": []*bot_common.BackgroundImageInfo{},
		"database_config":            []*bot_common.Database{},
		"shortcut_command":           []string{},
		"layout_info":                &bot_common.LayoutInfo{},
	},
	"single_agent_version": {
		"variable":                   []*bot_common.Variable{},
		"model_info":                 &bot_common.ModelInfo{},
		"onboarding_info":            &bot_common.OnboardingInfo{},
		"prompt":                     &bot_common.PromptInfo{},
		"plugin":                     []*bot_common.PluginInfo{},
		"knowledge":                  &bot_common.Knowledge{},
		"workflow":                   []*bot_common.WorkflowInfo{},
		"suggest_reply":              &bot_common.SuggestReplyInfo{},
		"jump_config":                &bot_common.JumpConfig{},
		"background_image_info_list": []*bot_common.BackgroundImageInfo{},
		"database_config":            []*bot_common.Database{},
		"shortcut_command":           []string{},
		"layout_info":                &bot_common.LayoutInfo{},
	},
	"plugin": {
		"manifest":    &pluginentity.PluginManifest{},
		"openapi_doc": &pluginentity.Openapi3T{},
		"ext":         map[string]any{},
	},
	"plugin_draft": {
		"manifest":    &pluginentity.PluginManifest{},
		"openapi_doc": &pluginentity.Openapi3T{},
	},
	"plugin_version": {
		"manifest":    &pluginentity.PluginManifest{},
		"openapi_doc": &pluginentity.Openapi3T{},
		"ext":         map[string]any{},
	},
	"agent_tool_draft": {
		"operation": &pluginentity.Openapi3Operation{},
	},
	"agent_tool_version": {
		"operation": &pluginentity.Openapi3Operation{},
	},
	"tool": {
		"operation": &pluginentity.Openapi3Operation{},
		"ext":       map[string]any{},
	},
	"tool_draft": {
		"operation": &pluginentity.Openapi3Operation{},
	},
	"tool_version": {
		"operation": &pluginentity.Openapi3Operation{},
		"ext":       map[string]any{},
	},
	"plugin_oauth_auth": {
		"oauth_config": &pluginentity.OAuthAuthorizationCodeConfig{},
	},
	"run_record": {
		"usage": &agentrunentity.Usage{},
	},
	"model_instance": {
		"provider":     &mgrentity.ModelProvider{},
		"display_info": &mgrentity.DisplayInfo{},
		"connection":   &mgrentity.Connection{},
		"capability":   &mgrentity.ModelAbility{},
		"parameters":   []mgrentity.ModelParameter{},
	},
}

func main() {
	db := initDB()
	generateForLLM(db)
	generateForFoundation(db)
	generateForComponent(db)
	generateForConversation(db)
	generateForData(db)
}

func initDB() *gorm.DB {
	var initSQLDir string
	if projectRoot, err := findProjectRoot(); err != nil {
		panic(err)
	} else {
		initSQLDir = filepath.Join(projectRoot, "deployment/bootstrap/mysql-init/init-sql")
	}
	cli, err := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{initSQLDir},
	}))
	if err != nil {
		panic(err)
	}
	return cli
}

func getGenerateConfig(path string) gen.Config {
	config := gen.Config{
		// 最终package不能设置为model，在有数据库表同步的情况下会产生冲突，若一定要使用可以单独指定model package的新名字
		OutPath:           fmt.Sprintf("./%s/query", path),
		ModelPkgPath:      fmt.Sprintf("./%s/model", path), // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		FieldNullable:     true,                            // 对于数据库中nullable的数据，在生成代码中自动对应为指针类型
		FieldWithIndexTag: true,                            // 从数据库同步的表结构代码包含gorm的index tag
		FieldWithTypeTag:  true,
	}
	return config
}

func generateForFoundation(db *gorm.DB) {
	var path string
	var tableList []string
	// User
	path = "modules/foundation/user/infra/repo/gorm_gen"
	tableList = []string{"user"}
	generateFunc(db, path, tableList)

	// OpenAuth
	path = "modules/foundation/openauth/infra/repo/gorm_gen"
	tableList = []string{"api_key"}
	generateFunc(db, path, tableList)

}

func generateForComponent(db *gorm.DB) {
	var path string
	var tableList []string

	path = "modules/component/agent/infra/repo/gorm_gen"
	tableList = []string{"single_agent_draft", "single_agent_publish", "single_agent_version"}
	generateFunc(db, path, tableList)

	path = "modules/component/prompt/infra/repo/gorm_gen"
	tableList = []string{"prompt_resource"}
	generateFunc(db, path, tableList)

	path = "modules/component/plugin/infra/repo/gorm_gen"
	tableList = []string{"agent_tool_draft", "agent_tool_version", "plugin",
		"plugin_draft", "plugin_oauth_auth", "plugin_version", "tool", "tool_draft", "tool_version"}
	generateFunc(db, path, tableList)
}

func generateForData(db *gorm.DB) {
	var path string
	var tableList []string

	path = "modules/data/upload/infra/repo/gorm_gen"
	tableList = []string{"files"}
	generateFunc(db, path, tableList)
}

func generateForConversation(db *gorm.DB) {
	var path string
	var tableList []string

	// Conversation
	path = "modules/conversation/conversation/infra/repo/gorm_gen"
	tableList = []string{"conversation"}
	generateFunc(db, path, tableList)

	// Message
	path = "modules/conversation/message/infra/repo/gorm_gen"
	tableList = []string{"message"}
	generateFunc(db, path, tableList)

	// Agent Run
	path = "modules/conversation/agent_run/infra/repo/gorm_gen"
	tableList = []string{"run_record"}
	generateFunc(db, path, tableList)
}

func generateForLLM(db *gorm.DB) {
	var path string
	var tableList []string

	path = "modules/llm/infra/repo/gorm_gen"
	tableList = []string{"model_entity", "model_meta", "model_instance", "model_request_record"}
	generateFunc(db, path, tableList)
}

func generateFunc(db *gorm.DB, projectPath string, tableList []string) {
	g := gen.NewGenerator(getGenerateConfig(projectPath))
	g.UseDB(db)

	var models []any
	for _, tableName := range tableList {
		var opts []gen.ModelOpt
		if col2Model, exist := table2Columns2Model[tableName]; exist {
			g.WithOpts(gen.FieldType("deleted_at", "gorm.DeletedAt"))
			genModify := func(col string, model any) func(f gen.Field) gen.Field {
				return func(f gen.Field) gen.Field {
					if f.ColumnName != col {
						return f
					}

					st := reflect.TypeOf(model)
					//f.Name = st.Name()
					f.Type = resolveType(st, true, projectPath)
					f.GORMTag.Set("serializer", "json")
					return f
				}
			}
			timeModify := func(f gen.Field) gen.Field {
				if f.ColumnName == "updated_at" {
					// https://gorm.io/zh_CN/docs/models.html#%E5%88%9B%E5%BB%BA-x2F-%E6%9B%B4%E6%96%B0%E6%97%B6%E9%97%B4%E8%BF%BD%E8%B8%AA%EF%BC%88%E7%BA%B3%E7%A7%92%E3%80%81%E6%AF%AB%E7%A7%92%E3%80%81%E7%A7%92%E3%80%81Time%EF%BC%89
					f.GORMTag.Set("autoUpdateTime", "milli")
				}
				if f.ColumnName == "created_at" {
					f.GORMTag.Set("autoCreateTime", "milli")
				}
				return f
			}
			for column, m := range col2Model {
				cp := m
				opts = append(opts, gen.FieldModify(genModify(column, cp)))
			}
			opts = append(opts, gen.FieldModify(timeModify))

		}
		models = append(models, g.GenerateModel(tableName, opts...))
	}

	g.ApplyBasic(models...)

	g.Execute()
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	backendDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	if _, err := os.Stat(filepath.Join(backendDir, "deployment")); os.IsNotExist(err) {
		return "", fmt.Errorf("could not find 'domain' directory in backend path: %s", backendDir)
	}

	return backendDir, nil
}

func resolveType(typ reflect.Type, required bool, modelPath string) string {
	switch typ.Kind() {
	case reflect.Ptr:
		return resolveType(typ.Elem(), false, modelPath)
	case reflect.Slice:
		return "[]" + resolveType(typ.Elem(), required, modelPath)
	default:
		prefix := "*"
		if required {
			prefix = ""
		}

		if strings.HasSuffix(typ.PkgPath(), modelPath) {
			return prefix + typ.Name()
		}

		return prefix + typ.String()
	}
}
