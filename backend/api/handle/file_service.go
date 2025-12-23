package handle

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

// GetFile 下载本地存储的文件
// @router /static/files/*filepath [GET]
func GetFile(c *gin.Context) {
	// 获取请求的文件路径
	filepathParam := c.Param("filepath")

	// 安全检查：确保文件路径在 local_storage 目录内
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")
	if localStoragePath == "" {
		localStoragePath = "./deployment/local_storage"
	}

	// 构建完整文件路径
	fullPath := filepath.Join(localStoragePath, filepathParam)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file not found",
			"path":  filepathParam,
		})
		return
	}

	// 检查是否是文件
	info, err := os.Stat(fullPath)
	if err != nil {
		logs.Error("Failed to stat file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "path is a directory, not a file",
		})
		return
	}

	// 返回文件
	c.File(fullPath)
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	backendDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	if _, err := os.Stat(filepath.Join(backendDir, "deployment")); os.IsNotExist(err) {
		return "", fmt.Errorf("could not find 'domain' directory in backend path: %s", backendDir)
	}

	return backendDir, nil
}
