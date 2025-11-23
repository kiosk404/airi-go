package logs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	CtxKeyLogID = "A_LOGID"
)

type RotateHook struct {
	Filename   string
	MaxSize    int64
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	suffix     string
	fileInfo   os.FileInfo
}

func NewRotateHook(filename string) *RotateHook {
	return &RotateHook{
		Filename:   filename,
		MaxSize:    100 * 1024 * 1024,
		MaxBackups: 3,
		MaxAge:     7,
		LocalTime:  false,
	}
}

func (hook *RotateHook) rotate() error {
	if hook.fileInfo != nil && hook.fileInfo.Size() < hook.MaxSize {
		return nil
	}

	err := hook.cleanUp()
	if err != nil {
		return err
	}

	fileName := hook.Filename + hook.suffix
	err = os.Rename(hook.Filename, fileName)
	if err != nil {
		return err
	}

	go hook.deleteOldFiles()

	return nil
}

func (hook *RotateHook) cleanUp() error {
	files, err := filepath.Glob(hook.Filename + ".*")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for len(files) >= hook.MaxBackups {
		err := os.Remove(files[0])
		if err != nil {
			return err
		}
		files = files[1:]
	}

	return nil
}

func (hook *RotateHook) deleteOldFiles() {
	if hook.MaxAge <= 0 {
		return
	}

	files, err := filepath.Glob(hook.Filename + ".*")
	if err != nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -hook.MaxAge)

	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			continue
		}

		if fi.ModTime().Before(cutoff) {
			os.Remove(file)
		}
	}
}

func (hook *RotateHook) Fire(entry *logrus.Entry) error {
	if hook.fileInfo == nil {
		fi, err := os.Stat(hook.Filename)
		if err != nil {
			return err
		}
		hook.fileInfo = fi
	}

	err := hook.rotate()
	if err != nil {
		return err
	}

	return nil
}

func (hook *RotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type FileHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.Writer.Write(line)
	return err
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type ConsoleHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.Writer.Write(line)
	return err
}

func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type Logger struct {
	*logrus.Logger
}

func NewLogger(filename string) (*Logger, error) {
	logger := logrus.New()

	file, err := createFile(filename)
	if err != nil {
		return nil, err
	}

	callerPrettyfier := func(frame *runtime.Frame) (function string, file string) {
		_, filename, line, ok := runtime.Caller(11)
		if !ok {
			return "", ""
		}
		if strings.Contains(filename, "pkg/logs/log.go") {
			_, filename, line, ok = runtime.Caller(12)
			if !ok {
				return "", ""
			}
		}

		relPath, err := filepath.Rel(getRootDir(), filename)
		if err != nil {
			return "", ""
		}

		function = fmt.Sprintf("%s:%d", relPath, line)

		return function, ""
	}

	// 创建控制台格式化器（带颜色）
	consoleFormatter := &logrus.TextFormatter{
		ForceColors:      true, // 强制颜色输出
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettyfier,
	}

	// 创建文件格式化器（不带颜色）
	fileFormatter := &logrus.TextFormatter{
		DisableColors:    true, // 禁用颜色输出
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettyfier,
	}

	// 创建控制台Hook
	consoleHook := &ConsoleHook{
		Writer:    os.Stdout,
		Formatter: consoleFormatter,
	}

	// 创建文件Hook
	fileHook := &FileHook{
		Writer:    file,
		Formatter: fileFormatter,
	}

	// 添加Hooks
	logger.AddHook(consoleHook)
	logger.AddHook(fileHook)

	// 禁用默认输出
	logger.SetOutput(io.Discard)

	// 添加字段来包含代码行号
	logger.SetReportCaller(true)

	rotateHook := NewRotateHook(filename)
	logger.AddHook(rotateHook)

	return &Logger{logger}, nil
}

func (l *Logger) GetLogID(ctx context.Context) string {
	logID, _ := ctx.Value(CtxKeyLogID).(string)
	return logID
}

func getRootDir() string {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return rootDir
}

func createFile(filename string) (*os.File, error) {
	dir := filepath.Dir(filename) // 获取目录路径

	// 判断是否包含目录路径
	if dir != "." && dir != ".." && dir != string(filepath.Separator) {
		err := os.MkdirAll(dir, os.ModePerm) // 创建目录，如果目录已存在则忽略错误
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // 创建文件
	if err != nil {
		return nil, err
	}

	return file, nil
}
