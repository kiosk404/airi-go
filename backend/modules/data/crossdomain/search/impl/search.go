package impl

import (
	"context"

	crosssearch "github.com/kiosk404/airi-go/backend/modules/data/crossdomain/search"
	"github.com/kiosk404/airi-go/backend/modules/data/crossdomain/search/model"
	"github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
)

var defaultSVC crosssearch.Search

type impl struct {
	DomainSVC crosssearch.Search
}

func (i impl) SearchResources(ctx context.Context, req *model.SearchResourcesRequest) (resp *model.SearchResourcesResponse, err error) {
	return i.DomainSVC.SearchResources(ctx, req)
}

func InitDomainService(c service.Search) crosssearch.Search {
	defaultSVC = &impl{
		DomainSVC: c,
	}

	return defaultSVC
}
