package logs

import (
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

type Logger struct {
	*logrus.Logger
}

type OutputHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	LogLevels []logrus.Level
}

func (hook *OutputHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *OutputHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func NewLogger(filename string) (*Logger, error) {
	logger := logrus.New()

	file, err := createFile(filename)
	if err != nil {
		return nil, err
	}

	fileFormatter := &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			_, filename, line, ok := runtime.Caller(9)
			if !ok {
				return "", ""
			}
			if strings.Contains(filename, "backend/pkg/logger/log.go") {
				_, filename, line, ok = runtime.Caller(10)
				if !ok {
					return "", ""
				}
			}

			relPath, err := filepath.Rel(getRootDir(), filename)
			if err != nil {
				return "", ""
			}

			filename = " "
			function = fmt.Sprintf("%s:%d", relPath, line)
			return function, filename
		},
	}

	consoleFormatter := &logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		CallerPrettyfier: fileFormatter.CallerPrettyfier, // 复用CallerPrettyfier
	}

	consoleHook := &OutputHook{
		Writer:    os.Stdout,
		Formatter: consoleFormatter,
		LogLevels: logrus.AllLevels,
	}

	fileHook := &OutputHook{
		Writer:    file,
		Formatter: fileFormatter,
		LogLevels: logrus.AllLevels,
	}

	logger.SetOutput(io.Discard) // 禁用默认输出
	logger.SetReportCaller(true)

	// 添加两个hook
	logger.AddHook(consoleHook)
	logger.AddHook(fileHook)

	rotateHook := NewRotateHook(filename)
	logger.AddHook(rotateHook)

	return &Logger{logger}, nil
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

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("%s:%d - %s", frame.File, frame.Line, frame.Function)
}
