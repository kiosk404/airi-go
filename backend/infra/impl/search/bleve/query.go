package bleve

import (
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/kiosk404/airi-go/backend/infra/contract/search"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func (b *bleveClient) query2BleveQuery(q *Query) query.Query {
	if q == nil {
		return bleve.NewMatchAllQuery()
	}

	var bleveQ query.Query

	switch q.Type {
	case search.QueryTypeEqual:
		// Term query - exact match
		tq := bleve.NewTermQuery(fmt.Sprint(q.KV.Value))
		tq.SetField(q.KV.Key)
		bleveQ = tq

	case search.QueryTypeMatch:
		// Match query - full text search
		mq := bleve.NewMatchQuery(fmt.Sprint(q.KV.Value))
		mq.SetField(q.KV.Key)
		bleveQ = mq

	case search.QueryTypeMultiMatch:
		// MultiMatch - search across multiple fields
		disjuncts := make([]query.Query, 0, len(q.MultiMatchQuery.Fields))
		for _, field := range q.MultiMatchQuery.Fields {
			mq := bleve.NewMatchQuery(q.MultiMatchQuery.Query)
			mq.SetField(field)
			disjuncts = append(disjuncts, mq)
		}
		if q.MultiMatchQuery.Operator == "and" {
			bleveQ = bleve.NewConjunctionQuery(disjuncts...)
		} else {
			bleveQ = bleve.NewDisjunctionQuery(disjuncts...)
		}

	case search.QueryTypeNotExists:
		// Not exists - documents where field is missing
		// Use boolean query with must_not for field existence
		boolQ := bleve.NewBooleanQuery()

		// Match all documents
		boolQ.AddMust(bleve.NewMatchAllQuery())

		// Exclude documents that have this field (by using a wildcard on the field)
		wq := bleve.NewWildcardQuery("*")
		wq.SetField(q.KV.Key)
		boolQ.AddMustNot(wq)

		bleveQ = boolQ

	case search.QueryTypeContains:
		// Wildcard/contains query
		wildcardStr := fmt.Sprintf("*%s*", q.KV.Value)
		wq := bleve.NewWildcardQuery(wildcardStr)
		wq.SetField(q.KV.Key)
		bleveQ = wq

	case search.QueryTypeIn:
		// Terms query - match any of the values
		disjuncts := make([]query.Query, 0)
		if values, ok := q.KV.Value.([]interface{}); ok {
			for _, val := range values {
				tq := bleve.NewTermQuery(fmt.Sprint(val))
				tq.SetField(q.KV.Key)
				disjuncts = append(disjuncts, tq)
			}
		}
		bleveQ = bleve.NewDisjunctionQuery(disjuncts...)

	default:
		bleveQ = bleve.NewMatchAllQuery()
	}

	// Handle Bool queries
	if q.Bool != nil {
		boolQuery := bleve.NewBooleanQuery()

		// Must clauses
		for _, must := range q.Bool.Must {
			boolQuery.AddMust(b.query2BleveQuery(&must))
		}

		// Should clauses
		for _, should := range q.Bool.Should {
			boolQuery.AddShould(b.query2BleveQuery(&should))
		}

		// MustNot clauses
		for _, mustNot := range q.Bool.MustNot {
			boolQuery.AddMustNot(b.query2BleveQuery(&mustNot))
		}

		// Filter clauses (treated as Must in Bleve)
		for _, filter := range q.Bool.Filter {
			boolQuery.AddMust(b.query2BleveQuery(&filter))
		}

		return boolQuery
	}

	return bleveQ
}

func (b *bleveClient) convertSort(sorts []search.SortFiled) []string {
	sortStrs := make([]string, 0, len(sorts))
	for _, s := range sorts {
		if s.Asc {
			sortStrs = append(sortStrs, s.Field)
		} else {
			sortStrs = append(sortStrs, "-"+s.Field)
		}
	}
	return sortStrs
}

func (b *bleveClient) convertSearchResult(result *bleve.SearchResult, req *Request) *Response {
	resp := &Response{
		Hits: search.HitsMetadata{
			Total: &search.TotalHits{
				Value: int64(result.Total),
			},
			MaxScore: ptr.Of(result.MaxScore),
			Hits:     make([]search.Hit, 0, len(result.Hits)),
		},
	}

	for _, hit := range result.Hits {
		h := search.Hit{
			Id_:     ptr.Of(hit.ID),
			Score_:  ptr.Of(hit.Score),
			Source_: json.RawMessage{},
		}

		resp.Hits.Hits = append(resp.Hits.Hits, h)
	}

	return resp
}

func (b *bleveClient) convertPropertyToFieldMapping(prop any) *mapping.FieldMapping {
	fieldMapping := bleve.NewTextFieldMapping()

	// Try to determine field type from property
	if propMap, ok := prop.(map[string]interface{}); ok {
		if typeVal, ok := propMap["type"].(string); ok {
			switch typeVal {
			case "long", "integer":
				fieldMapping = bleve.NewNumericFieldMapping()
			case "text":
				fieldMapping = bleve.NewTextFieldMapping()
			case "keyword":
				fieldMapping = bleve.NewKeywordFieldMapping()
			case "date":
				fieldMapping = bleve.NewDateTimeFieldMapping()
			case "boolean":
				fieldMapping = bleve.NewBooleanFieldMapping()
			}
		}
	}

	return fieldMapping
}
