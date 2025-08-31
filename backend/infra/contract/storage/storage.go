package storage

import (
	"context"
	"io"
	"time"
)

//go:generate  mockgen -destination ../../../internal/mock/infra/contract/storage/storage_mock.go -package mock -source storage.go Factory
type Storage interface {
	// PutObject puts the object with the specified key.
	PutObject(ctx context.Context, objectKey string, content []byte, opts ...PutOptFn) error
	// PutObjectWithReader puts the object with the specified key.
	PutObjectWithReader(ctx context.Context, objectKey string, content io.Reader, opts ...PutOptFn) error
	// GetObject returns the object with the specified key.
	GetObject(ctx context.Context, objectKey string) ([]byte, error)
	// DeleteObject deletes the object with the specified key.
	DeleteObject(ctx context.Context, objectKey string) error
	// GetObjectUrl returns a presigned URL for the object.
	// The URL is valid for the specified duration.
	GetObjectUrl(ctx context.Context, objectKey string, opts ...GetOptFn) (string, error)
	// ListAllObjects returns all objects with the specified prefix.
	// It may return a large number of objects, consider using ListObjectsPaginated for better performance.
	ListAllObjects(ctx context.Context, prefix string, withTagging bool) ([]*FileInfo, error)
	// ListObjectsPaginated returns objects with pagination support.
	// Use this method when dealing with large number of objects.
	ListObjectsPaginated(ctx context.Context, input *ListObjectsPaginatedInput) (*ListObjectsPaginatedOutput, error)
}

type SecurityToken struct {
	AccessKeyID     string `thrift:"access_key_id,1" frugal:"1,default,string" json:"access_key_id"`
	SecretAccessKey string `thrift:"secret_access_key,2" frugal:"2,default,string" json:"secret_access_key"`
	SessionToken    string `thrift:"session_token,3" frugal:"3,default,string" json:"session_token"`
	ExpiredTime     string `thrift:"expired_time,4" frugal:"4,default,string" json:"expired_time"`
	CurrentTime     string `thrift:"current_time,5" frugal:"5,default,string" json:"current_time"`
}

type ListObjectsPaginatedInput struct {
	Prefix   string
	PageSize int
	Cursor   string
	// Include objects tagging in the listing
	WithTagging bool
}

type ListObjectsPaginatedOutput struct {
	Files  []*FileInfo
	Cursor string
	// false: All results have been returned
	// true: There are more results to return
	IsTruncated bool
}
type FileInfo struct {
	Key          string
	LastModified time.Time
	ETag         string
	Size         int64
	Tagging      map[string]string
}
