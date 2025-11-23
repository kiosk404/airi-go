package bleve

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/bytedance/gg/gstd/gsync"
	"github.com/kiosk404/airi-go/backend/infra/contract/search"
)

type bleveClient struct {
	indexes      gsync.Map[string, bleve.Index]
	indexPathDir string
	types        *bleveTypes
}

type bleveBulkIndexer struct {
	index bleve.Index
	batch *bleve.Batch
}

func (b bleveBulkIndexer) Add(ctx context.Context, item search.BulkIndexerItem) error {
	switch item.Action {
	case "index", "create":
		return b.batch.Index(item.DocumentID, item.Body)
	case "delete":
		b.batch.Delete(item.DocumentID)
		return nil
	case "update":
		// For update, we need to index the document
		return b.batch.Index(item.DocumentID, item.Body)
	default:
		return fmt.Errorf("unsupported action: %s", item.Action)
	}
}

func (b bleveBulkIndexer) Close(ctx context.Context) error {
	return b.index.Batch(b.batch)
}

type bleveTypes struct{}

func (t *bleveTypes) NewLongNumberProperty() any {
	return map[string]interface{}{"type": "long"}
}

func (t *bleveTypes) NewTextProperty() any {
	return map[string]interface{}{"type": "text"}
}

func (t *bleveTypes) NewUnsignedLongNumberProperty() any {
	return map[string]interface{}{"type": "long"}
}

func (b *bleveClient) getIndex(idxName string) (bleve.Index, error) {
	var index bleve.Index
	var err error
	indexDir, indexPath := getEnvDefaultIndexPath(idxName)
	b.indexPathDir = indexDir
	if _, err = os.Stat(indexPath); os.IsNotExist(err) {
		// 不存在 → 创建新的 Index
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, err
		}
	} else {
		// 存在 → 打开已有 Index
		index, err = bleve.Open(indexPath)
		if err != nil {
			return nil, err
		}
	}
	idx, _ := b.indexes.LoadOrStore(idxName, index)
	return idx, nil
}

func newBleve() (Client, error) {

	return &bleveClient{indexes: gsync.Map[string, bleve.Index]{}}, nil
}

func (b *bleveClient) Create(ctx context.Context, index, id string, document any) error {
	idx, err := b.getIndex(index)
	if err != nil {
		return err
	}
	return idx.Index(id, document)
}

func (b *bleveClient) Update(ctx context.Context, index, id string, document any) error {
	idx, err := b.getIndex(index)
	if err != nil {
		return err
	}
	// 等价与Create
	return idx.Index(id, document)
}

func (b *bleveClient) Delete(ctx context.Context, index, id string) error {
	idx, err := b.getIndex(index)
	if err != nil {
		return err
	}
	// 等价与Create
	return idx.Delete(id)
}

func (b *bleveClient) Search(ctx context.Context, index string, req *Request) (*Response, error) {
	idx, err := b.getIndex(index)
	if err != nil {
		return nil, err
	}

	bleveQuery := b.query2BleveQuery(req.Query)

	searchReq := bleve.NewSearchRequest(bleveQuery)

	// Set size
	if req.Size != nil {
		searchReq.Size = *req.Size
	} else {
		searchReq.Size = 10 // default
	}

	// Set from
	if req.From != nil {
		searchReq.From = *req.From
	}

	// Set sorting
	if len(req.Sort) > 0 {
		searchReq.SortBy(b.convertSort(req.Sort))
	}

	// Set fields to return (important for getting document data)
	// Use "*" to return all fields
	searchReq.Fields = []string{"*"}

	searchResult, err := idx.Search(searchReq)
	if err != nil {
		return nil, err
	}

	// Convert Bleve result to ES-compatible Response
	return b.convertSearchResult(searchResult, req), nil
}

func (b *bleveClient) Exists(ctx context.Context, index string) (bool, error) {
	_, indexPath := getEnvDefaultIndexPath(index)
	_, err := os.Stat(indexPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (b *bleveClient) CreateIndex(ctx context.Context, index string, properties map[string]any) error {
	idxPath := filepath.Join(b.indexPathDir, index)

	// Create index mapping
	indexMapping := bleve.NewIndexMapping()

	// Convert ES properties to Bleve field mappings
	docMapping := bleve.NewDocumentMapping()
	for fieldName, prop := range properties {
		fieldMapping := b.convertPropertyToFieldMapping(prop)
		docMapping.AddFieldMappingsAt(fieldName, fieldMapping)
	}

	indexMapping.AddDocumentMapping("_default", docMapping)

	// Create the index
	idx, err := bleve.New(idxPath, indexMapping)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", index, err)
	}

	b.indexes.Store(index, idx)
	return nil
}

func (b *bleveClient) DeleteIndex(ctx context.Context, index string) error {
	idxPath := filepath.Join(b.indexPathDir, index)
	idx, exist := b.indexes.Load(idxPath)
	if exist && idx != nil {
		if err := idx.Close(); err != nil {
			return err
		}
	}
	b.indexes.Delete(idxPath)
	return nil
}

func (b *bleveClient) Types() search.Types {
	return &bleveTypes{}
}

func (b *bleveClient) NewBulkIndexer(index string) (search.BulkIndexer, error) {
	idx, err := b.getIndex(index)
	if err != nil {
		return nil, err
	}

	return &bleveBulkIndexer{
		index: idx,
		batch: idx.NewBatch(),
	}, nil
}
