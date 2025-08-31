package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"time"

	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/infra/impl/storage/proxy"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioClient struct {
	client          *minio.Client
	accessKeyID     string
	secretAccessKey string
	bucketName      string
	endpoint        string
}

func New(ctx context.Context, endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (storage.Storage, error) {
	m, err := getMinioClient(ctx, endpoint, accessKeyID, secretAccessKey, bucketName, useSSL)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func getMinioClient(_ context.Context, endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*minioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client failed %v", err)
	}

	m := &minioClient{
		client:          client,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		bucketName:      bucketName,
		endpoint:        endpoint,
	}

	err = m.createBucketIfNeed(context.Background(), client, bucketName, "cn-north-1")
	if err != nil {
		return nil, fmt.Errorf("init minio client failed %v", err)
	}

	// m.test()
	return m, nil
}

func (m *minioClient) createBucketIfNeed(ctx context.Context, client *minio.Client, bucketName, region string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("check bucket %s exist failed %v", bucketName, err)
	}

	if exists {
		return nil
	}

	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
	if err != nil {
		return fmt.Errorf("create bucket %s failed %v", bucketName, err)
	}

	return nil
}

func (m *minioClient) test() {
	ctx := context.Background()
	objectName := fmt.Sprintf("test-file-%d.txt", rand.Int())

	err := m.PutObject(ctx, objectName, []byte("hello content"),
		storage.WithContentType("text/plain"), storage.WithTagging(map[string]string{
			"uid":             "7543149965070155780",
			"conversation_id": "7543149965070155781",
			"type":            "user",
		}))
	if err != nil {
		logs.Error("upload file failed: %v", err)
	}

	logs.Info("upload file success")

	files, err := m.ListAllObjects(ctx, "test-file-", true)
	if err != nil {
		logs.Error("list objects failed: %v", err)
	}

	logs.Info("list objects success, files.len: %v", len(files))

	url, err := m.GetObjectUrl(ctx, objectName)
	if err != nil {
		logs.Error("get file url failed: %v", err)
	}

	logs.Info("get file url success, url: %s", url)

	content, err := m.GetObject(ctx, objectName)
	if err != nil {
		logs.Error("download file failed: %v", err)
	}

	logs.Info("download file success, content: %s", string(content))

	err = m.DeleteObject(ctx, objectName)
	if err != nil {
		logs.Error("delete object failed: %v", err)
	}

	logs.Info("delete object success")
}

func (m *minioClient) PutObject(ctx context.Context, objectKey string, content []byte, opts ...storage.PutOptFn) error {
	opts = append(opts, storage.WithObjectSize(int64(len(content))))
	return m.PutObjectWithReader(ctx, objectKey, bytes.NewReader(content), opts...)
}

func (m *minioClient) PutObjectWithReader(ctx context.Context, objectKey string, content io.Reader, opts ...storage.PutOptFn) error {
	option := storage.PutOption{}
	for _, opt := range opts {
		opt(&option)
	}

	minioOpts := minio.PutObjectOptions{}
	if option.ContentType != nil {
		minioOpts.ContentType = *option.ContentType
	}

	if option.ContentEncoding != nil {
		minioOpts.ContentEncoding = *option.ContentEncoding
	}

	if option.ContentDisposition != nil {
		minioOpts.ContentDisposition = *option.ContentDisposition
	}

	if option.ContentLanguage != nil {
		minioOpts.ContentLanguage = *option.ContentLanguage
	}

	if option.Expires != nil {
		minioOpts.Expires = *option.Expires
	}

	if option.Tagging != nil {
		minioOpts.UserTags = option.Tagging
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectKey,
		content, option.ObjectSize, minioOpts)
	if err != nil {
		return fmt.Errorf("PutObject failed: %v", err)
	}
	return nil
}

func (m *minioClient) GetObject(ctx context.Context, objectKey string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}
	defer obj.Close()
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("ReadObject failed: %v", err)
	}
	return data, nil
}

func (m *minioClient) DeleteObject(ctx context.Context, objectKey string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("DeleteObject failed: %v", err)
	}
	return nil
}

func (m *minioClient) GetObjectUrl(ctx context.Context, objectKey string, opts ...storage.GetOptFn) (string, error) {
	option := storage.GetOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.Expire == 0 {
		option.Expire = 3600 * 24 * 7
	}

	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, time.Duration(option.Expire)*time.Second, reqParams)
	if err != nil {
		return "", fmt.Errorf("GetObjectUrl failed: %v", err)
	}

	ok, proxyURL := proxy.CheckIfNeedReplaceHost(ctx, presignedURL.String())
	if ok {
		return proxyURL, nil
	}

	return presignedURL.String(), nil
}

func (m *minioClient) ListObjectsPaginated(ctx context.Context, input *storage.ListObjectsPaginatedInput) (*storage.ListObjectsPaginatedOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}
	if input.PageSize <= 0 {
		return nil, fmt.Errorf("page size must be positive")
	}

	files, err := m.ListAllObjects(ctx, input.Prefix, input.WithTagging)
	if err != nil {
		return nil, err
	}

	return &storage.ListObjectsPaginatedOutput{
		Files:       files,
		IsTruncated: false,
		Cursor:      "",
	}, nil
}

func (m *minioClient) ListAllObjects(ctx context.Context, prefix string, withTagging bool) ([]*storage.FileInfo, error) {
	opts := minio.ListObjectsOptions{
		Prefix:       prefix,
		Recursive:    true,
		WithMetadata: withTagging,
	}

	objectCh := m.client.ListObjects(ctx, m.bucketName, opts)

	var files []*storage.FileInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		files = append(files, &storage.FileInfo{
			Key:          object.Key,
			LastModified: object.LastModified,
			ETag:         object.ETag,
			Size:         object.Size,
			Tagging:      object.UserTags,
		})

		logs.Debug("key = %s, lastModified = %s, eTag = %s, size = %d, tagging = %v",
			object.Key, object.LastModified, object.ETag, object.Size, object.UserTags)
	}

	return files, nil
}
