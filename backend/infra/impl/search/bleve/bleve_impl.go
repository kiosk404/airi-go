package bleve

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/search"
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
	return newBleve()
}
