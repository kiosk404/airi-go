package logs

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	instance *Logger
	once     sync.Once
)

func Debug(format string, args ...interface{}) {
	if instance == nil {
		logrus.Debugf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.Debug(format)
	} else {
		instance.Debugf(format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if instance == nil {
		logrus.Infof(format, args...)
		return
	}
	if len(args) == 0 {
		instance.Info(format)
	} else {
		instance.Infof(format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if instance == nil {
		logrus.Warnf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.Warn(format)
	} else {
		instance.Warnf(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if instance == nil {
		logrus.Errorf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.Error(format)
	} else {
		instance.Errorf(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if instance == nil {
		logrus.Fatalf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.Fatal(format)
	} else {
		instance.Fatalf(format, args...)
	}
}

func DebugX(field string, format string, args ...interface{}) {
	if instance == nil {
		logrus.WithField("module", field).Debugf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.WithField("module", field).Debug(format)
	} else {
		instance.WithField("module", field).Debugf(format, args...)
	}
}

func InfoX(field string, format string, args ...interface{}) {
	if instance == nil {
		logrus.WithField("module", field).Infof(format, args...)
		return
	}
	if len(args) == 0 {
		instance.WithField("module", field).Info(format)
	} else {
		instance.WithField("module", field).Infof(format, args...)
	}
}

func WarnX(field string, format string, args ...interface{}) {
	if instance == nil {
		logrus.WithField("module", field).Warnf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.WithField("module", field).Warn(format)
	} else {
		instance.WithField("module", field).Warnf(format, args...)
	}
}

func ErrorX(field string, format string, args ...interface{}) {
	if instance == nil {
		logrus.WithField("module", field).Errorf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.WithField("module", field).Error(format)
	} else {
		instance.WithField("module", field).Errorf(format, args...)
	}
}

func FatalX(field string, format string, args ...interface{}) {
	if instance == nil {
		logrus.WithField("module", field).Fatalf(format, args...)
		return
	}
	if len(args) == 0 {
		instance.WithField("module", field).Fatal(format)
	} else {
		instance.WithField("module", field).Fatalf(format, args...)
	}
}

func InitLog(ctx context.Context, output string) (err error) {
	once.Do(func() {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
		if output == "" {
			output = "output/log/common.log"
		}
		instance, err = NewLogger(output)
	})
	return
}
