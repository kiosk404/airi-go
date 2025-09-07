package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func main() {
	db := initDB()
	generateForModelRequestRecord(db)
	generateForFoundationUser(db)
	generateForFoundationOpenAuth(db)
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
		OutPath: fmt.Sprintf("./%s/query", path),
		// Mode: gen.WithoutContext,
		ModelPkgPath:      fmt.Sprintf("./%s/model", path), // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		FieldNullable:     true,                            // 对于数据库中nullable的数据，在生成代码中自动对应为指针类型
		FieldWithIndexTag: true,                            // 从数据库同步的表结构代码包含gorm的index tag
		FieldWithTypeTag:  true,
	}
	config.WithImportPkgPath(fmt.Sprintf("github.com/kiosk404/airi-go/backend/%s/model", path))
	return config
}

func generateForModelRequestRecord(db *gorm.DB) {
	path := "modules/llm/infra/repo/gorm_gen"
	g := gen.NewGenerator(getGenerateConfig(path))
	g.UseDB(db)

	var models []any
	for _, table := range []string{
		"model_request_record",
	} {
		models = append(models, g.GenerateModel(table,
			// 添加软删除字段
			gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
			gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
				return tag.Set("column:deleted_at;not null;default:0;softDelete:milli")
			}),
			gen.FieldGORMTag("*", func(tag field.GormTag) field.GormTag {
				return tag.Set("charset=utf8mb4")
			})))
	}

	g.ApplyBasic(models...)
	g.Execute()
}

func generateForFoundationUser(db *gorm.DB) {
	path := "modules/foundation/user/infra/repo/gorm_gen"
	g := gen.NewGenerator(getGenerateConfig(path))
	g.UseDB(db)

	var models []any
	for _, table := range []string{
		"user",
	} {
		models = append(models, g.GenerateModel(table,
			// 添加软删除字段
			gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
			gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
				return tag.Set("column:deleted_at;not null;default:0;softDelete:milli")
			}),
			gen.FieldGORMTag("*", func(tag field.GormTag) field.GormTag {
				return tag.Set("charset=utf8mb4")
			})))
	}

	g.ApplyBasic(models...)
	g.Execute()
}

func generateForFoundationOpenAuth(db *gorm.DB) {
	path := "modules/foundation/openauth/infra/repo/gorm_gen"
	g := gen.NewGenerator(getGenerateConfig(path))
	g.UseDB(db)

	var models []any
	for _, table := range []string{
		"api_key",
	} {
		models = append(models, g.GenerateModel(table,
			// 添加软删除字段
			gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
			gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
				return tag.Set("column:deleted_at;not null;default:0;softDelete:milli")
			}),
			gen.FieldGORMTag("*", func(tag field.GormTag) field.GormTag {
				return tag.Set("charset=utf8mb4")
			})))
	}

	g.ApplyBasic(models...)
	g.Execute()
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	backendDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	if _, err := os.Stat(filepath.Join(backendDir, "deployment")); os.IsNotExist(err) {
		return "", fmt.Errorf("could not find 'domain' directory in backend path: %s", backendDir)
	}

	return backendDir, nil
}
