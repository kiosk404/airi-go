package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/kiosk404/airi-go/backend/domain/openauth/openapiauth/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var path2Table2Columns2Model = map[string]map[string]map[string]any{
	"domain/openauth/openapiauth/internal/dal/query": {
		"api_key": {},
	},
}

var fieldNullablePath = map[string]bool{}

const (
	// 存在内存中
	dbPathInMemory = ":memory:"
)

func main() {
	os.Setenv("LANG", "en_US.UTF-8")
	gormDB, err := gorm.Open(sqlite.Open(dbPathInMemory), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("gorm.Open failed, err=%v", err)
	}

	// 自动创建表结构
	for _, mapping := range path2Table2Columns2Model {
		for table := range mapping {
			model := getModelByTableName(table)
			if model != nil {
				// 创建表结构
				gormModel := convertToGormModel(model)
				err = gormDB.AutoMigrate(gormModel)
				if err != nil {
					log.Fatalf("failed to migrate table %s: %v", table, err)
				}
			}
		}
	}

	rootPath, err := findProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}
	fmt.Printf("rootPath: %s", rootPath)

	for path, mapping := range path2Table2Columns2Model {
		g := gen.NewGenerator(gen.Config{
			OutPath:       filepath.Join(rootPath, path),
			Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
			FieldNullable: fieldNullablePath[path],
		})

		parts := strings.Split(path, "/")
		modelPath := strings.Join(append(parts[:len(parts)-1], g.Config.ModelPkgPath), "/")

		g.UseDB(gormDB)
		g.WithOpts(gen.FieldType("deleted_at", "gorm.DeletedAt"))

		var resolveType func(typ reflect.Type, required bool) string
		resolveType = func(typ reflect.Type, required bool) string {
			switch typ.Kind() {
			case reflect.Ptr:
				return resolveType(typ.Elem(), false)
			case reflect.Slice:
				return "[]" + resolveType(typ.Elem(), required)
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

		genModify := func(col string, model any) func(f gen.Field) gen.Field {
			return func(f gen.Field) gen.Field {
				if f.ColumnName != col {
					return f
				}

				st := reflect.TypeOf(model)
				// f.Name = st.Name()
				f.Type = resolveType(st, true)
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

		var models []any
		for table, col2Model := range mapping {
			opts := make([]gen.ModelOpt, 0, len(col2Model))
			for column, m := range col2Model {
				cp := m
				opts = append(opts, gen.FieldModify(genModify(column, cp)))
			}
			opts = append(opts, gen.FieldModify(timeModify))
			models = append(models, g.GenerateModel(table, opts...))
		}

		g.ApplyBasic(models...)

		g.Execute()
	}
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	backendDir := filepath.Dir(filepath.Dir(filepath.Dir(filename))) // notice: the relative path of the script file is assumed here

	if _, err := os.Stat(filepath.Join(backendDir, "domain")); os.IsNotExist(err) {
		return "", fmt.Errorf("could not find 'domain' directory in backend path: %s", backendDir)
	}

	return backendDir, nil
}

// getModelByTableName 根据表名获取对应的模型结构体
func getModelByTableName(tableName string) interface{} {
	switch tableName {
	case "api_key":
		return &entity.ApiKey{}
	default:
		return nil
	}
}

// convertToGormModel 将实体模型转换为GORM模型(添加GORM标签)
func convertToGormModel(model interface{}) interface{} {
	// 这里可以根据需要添加GORM标签
	// 当前示例直接返回原模型，实际应用中可能需要添加标签
	return model
}
