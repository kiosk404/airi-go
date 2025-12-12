package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	localstorage "github.com/sajari/storage"
)

type LocalClient struct {
	store   localstorage.Local
	baseDir string
}

func New(ctx context.Context, pathDir string) (storage.Storage, error) {
	m, err := getLocalClient(pathDir)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getLocalClient(pathDir string) (*LocalClient, error) {
	// 确保目录存在
	if err := os.MkdirAll(pathDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", pathDir, err)
	}

	// 初始化一个基于本地文件系统的存储
	localStore := localstorage.Local(pathDir)

	c := LocalClient{
		store:   localStore,
		baseDir: pathDir,
	}

	return &c, nil
}

func (l *LocalClient) PutObject(ctx context.Context, objectKey string, content []byte, opts ...storage.PutOptFn) error {
	// 确保目录结构存在
	if err := l.ensureDir(objectKey); err != nil {
		return err
	}

	// 使用sajari/storage写入文件
	w, err := l.store.Create(ctx, objectKey)
	if err != nil {
		return fmt.Errorf("failed to create object %s: %w", objectKey, err)
	}
	defer w.Close()

	_, err = w.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write content to %s: %w", objectKey, err)
	}

	return nil
}

func (l *LocalClient) PutObjectWithReader(ctx context.Context, objectKey string, content io.Reader, opts ...storage.PutOptFn) error {
	// 确保目录结构存在
	if err := l.ensureDir(objectKey); err != nil {
		return err
	}

	// 使用sajari/storage写入文件
	w, err := l.store.Create(ctx, objectKey)
	if err != nil {
		return fmt.Errorf("failed to create object %s: %w", objectKey, err)
	}
	defer w.Close()

	_, err = io.Copy(w, content)
	if err != nil {
		return fmt.Errorf("failed to copy content to %s: %w", objectKey, err)
	}

	return nil
}

func (l *LocalClient) GetObject(ctx context.Context, objectKey string) ([]byte, error) {
	r, err := l.store.Open(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to open object %s: %w", objectKey, err)
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read content from %s: %w", objectKey, err)
	}

	return content, nil
}

func (l *LocalClient) DeleteObject(ctx context.Context, objectKey string) error {
	err := l.store.Delete(ctx, objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectKey, err)
	}
	return nil
}

func (l *LocalClient) GetObjectUrl(ctx context.Context, objectKey string, opts ...storage.GetOptFn) (string, error) {
	fullPath := filepath.Join(l.baseDir, objectKey)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", objectKey, err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("object %s does not exist", objectKey)
	}

	// 使用 filepath.ToSlash 确保路径分隔符是 /
	httpPath := filepath.ToSlash(objectKey)

	// 获取服务器地址
	host := "http://127.0.0.1:9527"
	if envHost := os.Getenv("SERVER_HOST"); envHost != "" {
		host = envHost
	}

	return host + "/static/files/" + httpPath, nil
}

func (l *LocalClient) ListAllObjects(ctx context.Context, prefix string, withTagging bool) ([]*storage.FileInfo, error) {
	var fileInfos []*storage.FileInfo

	// 构建搜索路径
	searchPath := filepath.Join(l.baseDir, prefix)
	baseDir := l.baseDir

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果路径不存在，跳过
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 获取相对路径作为objectKey
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		// 转换为统一的分隔符
		objectKey := filepath.ToSlash(relPath)

		// 检查是否匹配前缀
		if !strings.HasPrefix(objectKey, prefix) {
			return nil
		}

		fileInfo := &storage.FileInfo{
			Key:          objectKey,
			Size:         info.Size(),
			LastModified: info.ModTime(),
		}

		fileInfos = append(fileInfos, fileInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return fileInfos, nil
}

func (l *LocalClient) ListObjectsPaginated(ctx context.Context, input *storage.ListObjectsPaginatedInput) (*storage.ListObjectsPaginatedOutput, error) {
	// 获取所有匹配的文件
	allFiles, err := l.ListAllObjects(ctx, input.Prefix, input.WithTagging)
	if err != nil {
		return nil, err
	}

	// 处理分页
	startIndex := 0
	if input.Cursor != "" {
		// 使用文件路径作为cursor
		for i, file := range allFiles {
			if file.Key == input.Cursor {
				startIndex = i + 1
				break
			}
		}
	}

	// 计算结束索引
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 1000 // 默认页大小
	}

	endIndex := startIndex + pageSize
	if endIndex > len(allFiles) {
		endIndex = len(allFiles)
	}

	// 获取当前页的文件
	pageFiles := allFiles[startIndex:endIndex]

	// 准备输出
	output := &storage.ListObjectsPaginatedOutput{
		Files:       pageFiles,
		IsTruncated: endIndex < len(allFiles),
	}

	return output, nil
}

// ensureDir 确保对象键对应的目录结构存在
func (l *LocalClient) ensureDir(objectKey string) error {
	dir := filepath.Dir(filepath.Join(l.baseDir, objectKey))
	return os.MkdirAll(dir, 0755)
}
