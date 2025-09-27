package application

import (
	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
)

type ServiceComponents struct {
	IDGen     idgen.IDGenerator
	DB        rdb.Provider
	TosClient storage.Storage
	ImageX    imagex.ImageX

	SingleAgentDomainSVC singleagent.SingleAgent
}
