package bleve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kiosk404/airi-go/backend/types/consts"
)

func getEnvDefaultIndexPath(idxName string) (idxDir, idxPathName string) {
	indexPathPrefix := os.Getenv(consts.BleveIndexPath)
	if indexPathPrefix == "" {
		indexPathPrefix = "bleve_index" // 默认目录
	}
	indexPathName := fmt.Sprintf("%s%s", indexPathPrefix, idxName)
	indexPathDir := os.Getenv(consts.LocalStoragePath)
	return idxDir, filepath.Join(indexPathDir, indexPathName)
}
