package service

import (
	"context"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/infra/contract/search"
	"github.com/kiosk404/airi-go/backend/modules/data/crossdomain/search/model"
	searchEntity "github.com/kiosk404/airi-go/backend/modules/data/search/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/data/search/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

var searchInstance *searchImpl

func NewDomainService(ctx context.Context, e search.Client) Search {
	return &searchImpl{
		esClient: e,
	}
}

type searchImpl struct {
	esClient search.Client
}

type fieldName string

const (
	fieldOfOwnerID        = "owner_id"
	fieldOfID             = "id"
	fieldOfAPPID          = "app_id"
	fieldOfName           = "name"
	fieldOfNameRaw        = "name.raw"
	fieldOfHasPublished   = "has_published"
	fieldOfStatus         = "status"
	fieldOfType           = "type"
	fieldOfIsFav          = "is_fav"
	fieldOfIsRecentlyOpen = "is_recently_open"

	resTypeSearchAll = -1
)

func (s *searchImpl) SearchProjects(ctx context.Context, req *searchEntity.SearchProjectsRequest) (resp *searchEntity.SearchProjectsResponse, err error) {
	logs.DebugX(pkg.ModelName, "[SearchProjects] search : %s", conv.DebugJsonToStr(req))
	searchReq := &search.Request{
		Query: &search.Query{
			Bool: &search.BoolQuery{},
		},
	}

	if req.ProjectID != 0 { // precise search
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfID, conv.Int64ToStr(req.ProjectID)))
	}

	if req.Name != "" {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewContainsQuery(fieldOfNameRaw, req.Name))
	}

	if req.IsPublished {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfHasPublished, conv.BoolToInt(req.IsPublished)))
	}

	if req.IsFav {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfIsFav, conv.BoolToInt(req.IsFav)))
	}

	if req.IsRecentlyOpen {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfIsRecentlyOpen, conv.BoolToInt(req.IsRecentlyOpen)))
	}

	if req.OwnerID > 0 {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfOwnerID, conv.Int64ToStr(req.OwnerID)))
	}

	if len(req.Status) > 0 {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewInQuery(fieldOfStatus, req.Status))
	}

	if len(req.Types) > 0 {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewInQuery(fieldOfType, req.Types))
	}

	reqLimit := 100
	if req.Limit > 0 {
		reqLimit = int(req.Limit)
	}

	realLimit := reqLimit + 1
	searchReq.Size = &realLimit

	if req.OrderFiledName == "" {
		req.OrderFiledName = model.FieldOfUpdateTime
	}

	searchReq.Sort = []search.SortFiled{
		{
			Field: req.OrderFiledName,
			Asc:   req.OrderAsc,
		},
	}

	if req.Cursor != "" && req.Cursor != "0" {
		searchReq.SearchAfter = []any{
			fieldValueCaster(req.OrderFiledName, req.Cursor),
		}
	}

	result, err := s.esClient.Search(ctx, projectIndexName, searchReq)
	if err != nil {
		logs.DebugX(pkg.ModelName, "[Serarch.DO] err : %v", err)
		return nil, err
	}

	hits := result.Hits.Hits

	hasMore := func() bool {
		if len(hits) > reqLimit {
			return true
		}
		return false
	}()

	if hasMore {
		hits = hits[:reqLimit]
	}

	docs := make([]*searchEntity.ProjectDocument, 0, len(hits))
	for _, hit := range hits {
		doc, err := hit2AppDocument(hit)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	nextCursor := ""
	if len(docs) > 0 {
		nextCursor = formatProjectNextCursor(req.OrderFiledName, docs[len(docs)-1])
	}
	if nextCursor == "" {
		hasMore = false
	}

	resp = &searchEntity.SearchProjectsResponse{
		Data:       docs,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}

	return resp, nil
}

func hit2AppDocument(hit search.Hit) (*searchEntity.ProjectDocument, error) {
	doc := &searchEntity.ProjectDocument{}

	if err := sonic.Unmarshal(hit.Source_, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func fieldValueCaster(fieldName, cursor string) any {
	switch fieldName {
	case model.FieldOfCreateTime,
		model.FieldOfUpdateTime,
		model.FieldOfPublishTime,
		model.FieldOfFavTime,
		model.FieldOfRecentlyOpenTime:
		cursorInt, err := strconv.ParseInt(cursor, 10, 64)
		if err != nil {
			cursorInt = 0
		}

		return cursorInt
	default:
		return cursor
	}
}

func formatProjectNextCursor(ob string, val *searchEntity.ProjectDocument) string {
	fieldName2Cursor := map[string]string{
		model.FieldOfCreateTime:       conv.Int64ToStr(val.GetCreateTime()),
		model.FieldOfUpdateTime:       conv.Int64ToStr(val.GetUpdateTime()),
		model.FieldOfPublishTime:      conv.Int64ToStr(val.GetPublishTime()),
		model.FieldOfFavTime:          conv.Int64ToStr(val.GetFavTime()),
		model.FieldOfRecentlyOpenTime: conv.Int64ToStr(val.GetRecentlyOpenTime()),
	}

	res, ok := fieldName2Cursor[ob]
	if !ok {
		return ""
	}

	return res
}

func formatResourceNextCursor(ob string, val *searchEntity.ResourceDocument) string {
	fieldName2Cursor := map[string]string{
		model.FieldOfCreateTime:  conv.Int64ToStr(val.GetCreateTime()),
		model.FieldOfUpdateTime:  conv.Int64ToStr(val.GetUpdateTime()),
		model.FieldOfPublishTime: conv.Int64ToStr(val.GetPublishTime()),
	}

	res, ok := fieldName2Cursor[ob]
	if !ok {
		return ""
	}

	return res
}

func (s *searchImpl) SearchResources(ctx context.Context, req *searchEntity.SearchResourcesRequest) (resp *searchEntity.SearchResourcesResponse, err error) {
	logs.DebugX(pkg.ModelName, "[SearchResources] search : %s", conv.DebugJsonToStr(req))
	searchReq := &search.Request{
		Query: &search.Query{
			Bool: &search.BoolQuery{},
		},
	}

	if req.APPID > 0 {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfAPPID, conv.Int64ToStr(req.APPID)))
	} else {
		searchReq.Query.Bool.Should = append(searchReq.Query.Bool.Should,
			search.NewNotExistsQuery(fieldOfAPPID))
		searchReq.Query.Bool.Should = append(searchReq.Query.Bool.Should,
			search.NewEqualQuery(fieldOfAPPID, "0"))

		searchReq.Query.Bool.MinimumShouldMatch = ptr.Of(1)
	}

	if req.Name != "" {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewContainsQuery(fieldOfNameRaw, req.Name))
	}

	if req.OwnerID > 0 {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(fieldOfOwnerID, conv.Int64ToStr(req.OwnerID)))
	}

	if len(req.ResTypeFilter) == 1 && int(req.ResTypeFilter[0]) != resTypeSearchAll {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(searchEntity.FieldOfResType, req.ResTypeFilter[0]))
	}

	if len(req.ResTypeFilter) == 2 {
		resType := req.ResTypeFilter[0]
		resSubType := int(req.ResTypeFilter[1])

		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(searchEntity.FieldOfResType, resType))

		if resSubType != resTypeSearchAll {
			searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
				search.NewEqualQuery(searchEntity.FieldOfResSubType, int(resSubType)))
		}
	}

	if req.PublishStatusFilter != 0 {
		searchReq.Query.Bool.Must = append(searchReq.Query.Bool.Must,
			search.NewEqualQuery(searchEntity.FieldOfPublishStatus, req.PublishStatusFilter))
	}

	if req.OrderFiledName == "" {
		req.OrderFiledName = model.FieldOfUpdateTime
	}

	searchReq.Sort = []search.SortFiled{
		{
			Field: req.OrderFiledName,
			Asc:   req.OrderAsc,
		},
	}

	reqLimit := 100
	if req.Limit > 0 {
		reqLimit = int(req.Limit)
	}
	realLimit := reqLimit + 1
	searchReq.Size = &realLimit

	if req.Page != nil {
		page := *req.Page
		if page <= 0 {
			page = 1
		}
		searchReq.From = ptr.Of(int(page-1) * reqLimit)
	} else if req.Cursor != "" && req.Cursor != "0" {
		searchReq.SearchAfter = []any{
			req.Cursor,
		}
	}

	result, err := s.esClient.Search(ctx, resourceIndexName, searchReq)
	if err != nil {
		return nil, err
	}

	hits := result.Hits.Hits

	hasMore := func() bool {
		if len(hits) > reqLimit {
			return true
		}
		return false
	}()

	if hasMore {
		hits = hits[:reqLimit]
	}

	docs := make([]*searchEntity.ResourceDocument, 0, len(hits))
	for _, hit := range hits {
		doc := &searchEntity.ResourceDocument{}
		if err = sonic.Unmarshal(hit.Source_, doc); err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	nextCursor := ""
	if len(docs) > 0 {
		nextCursor = formatResourceNextCursor(req.OrderFiledName, docs[len(docs)-1])
	}
	if nextCursor == "" {
		hasMore = false
	}

	var total *int64
	if result.Hits.Total != nil {
		total = ptr.Of(result.Hits.Total.Value)
	}

	resp = &searchEntity.SearchResourcesResponse{
		Data:       docs,
		TotalHits:  total,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}

	return resp, nil
}
