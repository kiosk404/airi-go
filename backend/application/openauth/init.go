package openauth

import (
	"github.com/kiosk404/airi-go/backend/domain/openauth/openapiauth"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"gorm.io/gorm"
)

var (
	openapiAuthDomainSVC openapiauth.APIAuth
)

func InitService(db *gorm.DB, idGenSVC idgen.IDGenerator) *OpenAuthApplicationService {
	openapiAuthDomainSVC = openapiauth.NewService(&openapiauth.Components{
		IDGen: idGenSVC,
		DB:    db,
	})

	OpenAuthApplication.OpenAPIDomainSVC = openapiAuthDomainSVC

	return OpenAuthApplication
}
