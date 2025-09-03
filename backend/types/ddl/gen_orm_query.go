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

var path2Table2Model = map[string]map[string]any{
	"domain/openauth/openapiauth/internal/dal/query": {
		"api_key": &entity.ApiKey{},
	},
}

var fieldNullablePath = map[string]bool{}

const (
	// 存在内存中
	dbPathInMemory = ":memory:"

	// 存在 sqlite 中
	dbPathInSqlite = "airi-go.db"
)

var sqliteDBPath = func() string {
	if rootP, err := findProjectRoot(); err != nil {
		return ""
	} else {
		return fmt.Sprintf("%s/%s", rootP, dbPathInSqlite)
	}
}()

func main() {
	_ = os.Setenv("LANG", "en_US.UTF-8")
	gormDB, err := gorm.Open(sqlite.Open(dbPathInMemory), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("gorm.Open failed, err=%v", err)
	}

	// 自动创建表结构
	for _, mapping := range path2Table2Model {
		for table, model := range mapping {
			if err := gormDB.AutoMigrate(model); err != nil {
				log.Fatalf("AutoMigrate table %s failed, err=%v", table, err)
			}
		}
	}

	rootPath, err := findProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}
	fmt.Printf("rootPath: %s\n", rootPath)

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

		// ----------- 生成逻辑 ----------
		var models []any
		for table, col2Model := range mapping {
			opts := make([]gen.ModelOpt, 0, len(col2Model)+2) // +2 for timeModify and entityTypeModify

			// 保留原有的列特定修饰器逻辑
			for column, m := range col2Model {
				cp := m
				opts = append(opts, gen.FieldModify(genModify(column, cp)))
			}

			// 保留原有的时间修饰器
			opts = append(opts, gen.FieldModify(timeModify))

			// 新增：基于entity的类型修饰器（放在最后，确保覆盖前面的类型设置）
			if tableModels, exists := path2Table2Model[path]; exists {
				if entityModel, exists := tableModels[table]; exists && entityModel != nil {
					entityFieldTypes := extractEntityFieldTypes(entityModel)
					entityTypeModifier := createEntityTypeModifier(entityFieldTypes)
					opts = append(opts, gen.FieldModify(entityTypeModifier))
				}
			}

			models = append(models, g.GenerateModel(table, opts...))
		}

		if len(models) > 0 {
			g.ApplyBasic(models...)
		}

		g.Execute()
	}
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	backendDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	if _, err := os.Stat(filepath.Join(backendDir, "domain")); os.IsNotExist(err) {
		return "", fmt.Errorf("could not find 'domain' directory in backend path: %s", backendDir)
	}

	return backendDir, nil
}

// 从entity提取字段类型信息
func extractEntityFieldTypes(entityModel interface{}) map[string]string {
	fieldTypes := make(map[string]string)

	entityType := reflect.TypeOf(entityModel)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		// 获取json tag作为列名，如果没有则使用字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = strings.ToLower(field.Name)
		} else {
			// 处理json tag中的选项 (如 `json:"id,omitempty"`)
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonTag = jsonTag[:idx]
			}
		}

		// 直接使用json tag作为列名
		columnName := jsonTag

		// 获取Go类型
		fieldType := resolveUnderlyingType(field.Type)

		fieldTypes[columnName] = fieldType

		fmt.Printf("Entity field mapping: %s -> %s (type: %s)\n",
			field.Name, columnName, fieldType)
	}

	return fieldTypes
}

// 新增：基于entity的类型修饰器
func createEntityTypeModifier(entityFieldTypes map[string]string) func(gen.Field) gen.Field {
	return func(f gen.Field) gen.Field {
		// 如果entity中定义了这个字段的类型，就使用entity的类型
		if entityType, exists := entityFieldTypes[f.ColumnName]; exists {
			fmt.Printf("Applying entity type for %s: %s -> %s\n",
				f.ColumnName, f.Type, entityType)
			f.Type = entityType
		}
		return f
	}
}

// 解析底层类型，将自定义类型转换为其底层的基础类型
func resolveUnderlyingType(t reflect.Type) string {
	// 如果是指针类型，先获取元素类型
	if t.Kind() == reflect.Ptr {
		return "*" + resolveUnderlyingType(t.Elem())
	}

	// 如果是slice类型
	if t.Kind() == reflect.Slice {
		return "[]" + resolveUnderlyingType(t.Elem())
	}

	// 对于自定义类型，获取其底层类型
	if t.PkgPath() != "" && t.Kind() != reflect.Struct && t.Kind() != reflect.Interface {
		// 这是一个自定义类型（如 type AkType int32）
		underlyingType := t.Kind()
		switch underlyingType {
		case reflect.Int:
			return "int"
		case reflect.Int8:
			return "int8"
		case reflect.Int16:
			return "int16"
		case reflect.Int32:
			return "int32"
		case reflect.Int64:
			return "int64"
		case reflect.Uint:
			return "uint"
		case reflect.Uint8:
			return "uint8"
		case reflect.Uint16:
			return "uint16"
		case reflect.Uint32:
			return "uint32"
		case reflect.Uint64:
			return "uint64"
		case reflect.Float32:
			return "float32"
		case reflect.Float64:
			return "float64"
		case reflect.String:
			return "string"
		case reflect.Bool:
			return "bool"
		default:
			// 如果无法识别，返回原始类型名
			return t.String()
		}
	}

	// 对于基础类型或结构体，直接返回类型名
	return t.String()
}
