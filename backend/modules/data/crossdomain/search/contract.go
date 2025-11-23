package search

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/data/crossdomain/search/model"
)

type Search interface {
	SearchResources(ctx context.Context, req *model.SearchResourcesRequest) (resp *model.SearchResourcesResponse, err error)
}

var defaultSVC Search

func DefaultSVC() Search {
	return defaultSVC
}

func SetDefaultSVC(svc Search) {
	defaultSVC = svc
}
