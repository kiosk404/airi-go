package appinfra

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/impl/cache/local"
	idgenimpl "github.com/kiosk404/airi-go/backend/infra/impl/idgen"
	"github.com/kiosk404/airi-go/backend/infra/impl/rdb/mysql"
	"github.com/kiosk404/airi-go/backend/infra/impl/storage"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
)

type AppDependencies struct {
	DB            rdb.Provider
	CacheCli      cache.Cmdable
	IDGenSVC      idgen.IDGenerator
	TOSClient     storage.Storage
	ConfigFactory conf.IConfigLoaderFactory
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

func getMysqlDomain() string {
	return os.Getenv("AIRI_GO_MYSQL_DOMAIN")
}

func getMysqlPort() string {
	return os.Getenv("AIRI_GO_MYSQL_PORT")
}

func getMysqlUser() string {
	return os.Getenv("AIRI_GO_MYSQL_USER")
}

func getMysqlPassword() string {
	return os.Getenv("AIRI_GO_MYSQL_PASSWORD")
}

func getMysqlDatabase() string {
	return os.Getenv("AIRI_GO_MYSQL_DATABASE")
}
