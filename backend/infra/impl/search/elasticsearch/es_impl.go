package elasticsearch

import (
	"fmt"
	"os"

	"github.com/kiosk404/airi-go/backend/infra/contract/search"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

type (
	Client          = search.Client
	Types           = search.Types
	BulkIndexer     = search.BulkIndexer
	BulkIndexerItem = search.BulkIndexerItem
	BoolQuery       = search.BoolQuery
	Query           = search.Query
	Response        = search.Response
	Request         = search.Request
)

func New() (Client, error) {
	v := os.Getenv(consts.SearchESVersion)
	if v == "v8" {
		return newES8()
	} else if v == "v7" {
		return newES7()
	}

	return nil, fmt.Errorf("unsupported elasticsearch version %s", v)
}
