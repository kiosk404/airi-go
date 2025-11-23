package service

import (
	"context"
	"strings"

	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/service/official"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

type promptService struct {
	Repo repo.PromptRepository
}

func NewService(repo repo.PromptRepository) Prompt {
	return &promptService{
		Repo: repo,
	}
}

func (s *promptService) CreatePromptResource(ctx context.Context, p *entity.PromptResource) (int64, error) {
	return s.Repo.CreatePromptResource(ctx, p)
}

func (s *promptService) UpdatePromptResource(ctx context.Context, promptID int64, name, description, promptText *string) error {
	return s.Repo.UpdatePromptResource(ctx, promptID, name, description, promptText)
}

func (s *promptService) GetPromptResource(ctx context.Context, promptID int64) (*entity.PromptResource, error) {
	return s.Repo.GetPromptResource(ctx, promptID)
}

func (s *promptService) DeletePromptResource(ctx context.Context, promptID int64) error {
	err := s.Repo.DeletePromptResource(ctx, promptID)
	if err != nil {
		return err
	}

	return nil
}

func (s *promptService) ListOfficialPromptResource(ctx context.Context, keyword string) ([]*entity.PromptResource, error) {
	promptList := official.GetPromptList()

	promptList = searchPromptResourceList(ctx, promptList, keyword)
	return deepCopyPromptResource(promptList), nil
}

func deepCopyPromptResource(pl []*entity.PromptResource) []*entity.PromptResource {
	return slices.Transform(pl, func(p *entity.PromptResource) *entity.PromptResource {
		return &entity.PromptResource{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			PromptText:  p.PromptText,
			Status:      1,
		}
	})
}

func searchPromptResourceList(ctx context.Context, resource []*entity.PromptResource, keyword string) []*entity.PromptResource {
	if len(keyword) == 0 {
		return resource
	}

	retVal := make([]*entity.PromptResource, 0, len(resource))
	for _, promptResource := range resource {
		if promptResource == nil {
			continue
		}
		// name match
		if strings.Contains(strings.ToLower(promptResource.Name), strings.ToLower(keyword)) {
			retVal = append(retVal, promptResource)
			continue
		}
		// Body Match
		if strings.Contains(strings.ToLower(promptResource.PromptText), strings.ToLower(keyword)) {
			retVal = append(retVal, promptResource)
		}
	}
	return retVal
}
