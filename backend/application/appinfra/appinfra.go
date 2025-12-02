package appinfra

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kiosk404/airi-go/backend/api/model/llm/manage"
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/impl/cache/local"
	idgenimpl "github.com/kiosk404/airi-go/backend/infra/impl/idgen"
	"github.com/kiosk404/airi-go/backend/infra/impl/rdb/mysql"
	"github.com/kiosk404/airi-go/backend/infra/impl/storage"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

type AppDependencies struct {
	DB            rdb.Provider
	CacheCli      cache.Cmdable
	IDGenSVC      idgen.IDGenerator
	TOSClient     storage.Storage
	ImageXClient  imagex.ImageX
	ConfigFactory conf.IConfigLoaderFactory
	ModelMgr      manage.LLMManageService
}

func Init(ctx context.Context) (*AppDependencies, error) {
	deps := &AppDependencies{}
	var err error
	if deps.DB, err = mysql.NewDB(mysqlDBConfig()); err != nil {
		return nil, fmt.Errorf("init db failed, err=%w", err)
	}
	if deps.CacheCli, err = local.New(); err != nil {
		return nil, fmt.Errorf("init cache failed, err=%w", err)
	}
	if deps.IDGenSVC, err = idgenimpl.New(deps.CacheCli); err != nil {
		return nil, fmt.Errorf("init id gen svc failed, err=%w", err)
	}
	if deps.TOSClient, err = storage.New(ctx); err != nil {
		return nil, fmt.Errorf("init storage failed, err=%w", err)
	}
	if deps.ImageXClient, err = storage.NewImageX(ctx); err != nil {
		return nil, fmt.Errorf("init imagex client failed, err=%w", err)
	}
	if deps.ConfigFactory, err = modelmgr.ModelMetaConfFactory(getApplicationProjectRoot()); err != nil {
		return nil, fmt.Errorf("init model meta conf factory failed, err=%w", err)
	}

	return deps, err
}

func mysqlDBConfig() *mysql.Config {
	return &mysql.Config{
		DBHostname:   getMysqlDomain(),
		DBPort:       getMysqlPort(),
		User:         getMysqlUser(),
		Password:     getMysqlPassword(),
		DBName:       getMysqlDatabase(),
		Loc:          "Local",
		DBCharset:    "utf8mb4",
		Timeout:      time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		DSNParams:    url.Values{"clientFoundRows": []string{"true"}},
	}
}

func getApplicationProjectRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	_, err := os.Stat(filepath.Join(projectRoot, "conf"))
	if err != nil {
		return "."
	}
	return projectRoot
}

func getMysqlDomain() string {
	return os.Getenv(consts.MySQLDomain)
}

func getMysqlPort() string {
	return os.Getenv(consts.MySQLPort)
}

func getMysqlUser() string {
	return os.Getenv(consts.MySQLUser)
}

func getMysqlPassword() string {
	return os.Getenv(consts.MySQLPassport)
}

func getMysqlDatabase() string {
	return os.Getenv(consts.MySQLDatabase)
}
