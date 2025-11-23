package application

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/repo"
	prompt "github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/service"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	"gorm.io/gorm"
)

func InitService(db *gorm.DB, idGenSVC idgen.IDGenerator, re search.ResourceEventBus) *PromptApplicationService {
	repo := repo.NewPromptRepo(db, idGenSVC)
	PromptSVC.DomainSVC = prompt.NewService(repo)
	PromptSVC.eventbus = re

	return PromptSVC
}
