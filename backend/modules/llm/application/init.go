package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/repo"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
)

type (
	ModelManager = modelmgr.ModelManager
)

var ModelMgrSVC *ModelManagerApplicationService

type ServiceComponents struct {
	IDGen     idgen.IDGenerator
	DB        rdb.Provider
	TosClient storage.Storage
	IConf     conf.IConfigLoaderFactory
}

func InitService(ctx context.Context, idGen idgen.IDGenerator, db rdb.Provider,
	tosClient storage.Storage, iConf conf.IConfigLoaderFactory) *ModelManagerApplicationService {
	modelManageRepo := repo.NewModelMgrRepo(db, idGen)
	modelMgrDomainSVC := modelmgr.NewService(tosClient, modelManageRepo, iConf)
	ModelMgrSVC = newApplicationService(&ServiceComponents{
		IDGen:     idGen,
		DB:        db,
		TosClient: tosClient,
		IConf:     iConf,
	}, modelMgrDomainSVC)
	return ModelMgrSVC
}
