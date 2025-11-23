package search

import (
	"encoding/json"
	"io"
)

type BulkIndexerItem struct {
	Index           string
	Action          string
	DocumentID      string
	Routing         string
	Version         *int64
	VersionType     string
	Body            io.ReadSeeker
	RetryOnConflict *int
}

type Request struct {
	Size        *int
	Query       *Query
	MinScore    *float64
	Sort        []SortFiled
	SearchAfter []any
	From        *int
}

type SortFiled struct {
	Field string
	Asc   bool
}

type Response struct {
	Hits     HitsMetadata `json:"hits"`
	MaxScore *float64     `json:"max_score,omitempty"`
}

type HitsMetadata struct {
	Hits     []Hit    `json:"hits"`
	MaxScore *float64 `json:"max_score,omitempty"`
	// Total Total hit count information, present only if `track_total_hits` wasn't
	// `false` in the search request.
	Total *TotalHits `json:"total,omitempty"`
}

type Hit struct {
	Id_     *string         `json:"_id,omitempty"`
	Score_  *float64        `json:"_score,omitempty"`
	Source_ json.RawMessage `json:"_source,omitempty"`
}

type TotalHits struct {
	Relation TotalHitsRelation `json:"relation"`
	Value    int64             `json:"value"`
}

type TotalHitsRelation struct {
	Name string
}
