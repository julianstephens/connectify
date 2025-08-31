package config

import (
	"log"
	"sync"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/julianstephens/connectify/server/internal/utils"
)

var (
	logger     LoggerInterface
	loggerInit sync.Once
)

type LoggerInterface interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
}

type logrusWrapper struct {
	*logrus.Logger
}

func (l *logrusWrapper) Info(args ...interface{})  { l.Logger.Info(args...) }
func (l *logrusWrapper) Error(args ...interface{}) { l.Logger.Error(args...) }
func (l *logrusWrapper) Warn(args ...interface{})  { l.Logger.Warn(args...) }
func (l *logrusWrapper) Debug(args ...interface{}) { l.Logger.Debug(args...) }

type zapWrapper struct {
	*zap.Logger
}

func (z *zapWrapper) Info(args ...interface{})  { z.Logger.Sugar().Info(args...) }
func (z *zapWrapper) Error(args ...interface{}) { z.Logger.Sugar().Error(args...) }
func (z *zapWrapper) Warn(args ...interface{})  { z.Logger.Sugar().Warn(args...) }
func (z *zapWrapper) Debug(args ...interface{}) { z.Logger.Sugar().Debug(args...) }

// SetLogger sets the singleton logger based on LogType
func SetLogger(logType string) {
	loggerInit.Do(func() {
		switch logType {
		case "logrus":
			logger = &logrusWrapper{logrus.New()}
		default:
			zl, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("failed to initialize zap logger: %v", err)
			}
			logger = &zapWrapper{zl}
		}
	})
}

// GetLogger returns the singleton logger
func GetLogger() LoggerInterface {
	if logger == nil {
		SetLogger(utils.DefaultIfNil(utils.PointerTo(AppConfig.LogType), "zap"))
	}
	return logger
}
