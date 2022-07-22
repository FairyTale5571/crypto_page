package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	DebugCtx(ctx context.Context, args ...interface{})
	InfoCtx(ctx context.Context, args ...interface{})
	WarnCtx(ctx context.Context, args ...interface{})
	ErrorCtx(ctx context.Context, args ...interface{})
	FatalCtx(ctx context.Context, args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	DebugfCtx(ctx context.Context, format string, args ...interface{})
	InfofCtx(ctx context.Context, format string, args ...interface{})
	WarnfCtx(ctx context.Context, format string, args ...interface{})
	ErrorfCtx(ctx context.Context, format string, args ...interface{})
	FatalfCtx(ctx context.Context, format string, args ...interface{})
}

type LoggerWrapper struct {
	lg    *logrus.Logger
	entry *logrus.Entry
}

func New(service string) *LoggerWrapper {
	log := &LoggerWrapper{lg: logrus.New()}

	log.lg.SetFormatter(&logrus.JSONFormatter{})
	log.lg.SetOutput(os.Stdout)
	log.lg.SetLevel(logrus.DebugLevel)
	log.entry = log.lg.WithFields(logrus.Fields{
		"service": service,
	})
	return log
}

func (logger *LoggerWrapper) Debug(args ...interface{}) {
	logger.entry.Debug(args...)
}

func (logger *LoggerWrapper) Info(args ...interface{}) {
	logger.entry.Info(args...)
}

func (logger *LoggerWrapper) Warn(args ...interface{}) {
	logger.entry.Warn(args...)
}

func (logger *LoggerWrapper) Error(args ...interface{}) {
	logger.entry.Error(args...)
}

func (logger *LoggerWrapper) Fatal(args ...interface{}) {
	logger.entry.Fatal(args...)
}

func (logger *LoggerWrapper) DebugCtx(ctx context.Context, args ...interface{}) {
	logger.entry.Debug(args...)
}

func (logger *LoggerWrapper) InfoCtx(ctx context.Context, args ...interface{}) {
	logger.entry.Info(args...)
}

func (logger *LoggerWrapper) WarnCtx(ctx context.Context, args ...interface{}) {
	logger.entry.Warn(args...)
}

func (logger *LoggerWrapper) ErrorCtx(ctx context.Context, args ...interface{}) {
	logger.entry.Error(args...)
}

func (logger *LoggerWrapper) FatalCtx(ctx context.Context, args ...interface{}) {
	logger.entry.Fatal(args...)
}

func (logger *LoggerWrapper) Debugf(format string, args ...interface{}) {
	logger.entry.Debugf(format, args...)
}

func (logger *LoggerWrapper) Infof(format string, args ...interface{}) {
	logger.entry.Infof(format, args...)
}

func (logger *LoggerWrapper) Warnf(format string, args ...interface{}) {
	logger.entry.Warnf(format, args...)
}

func (logger *LoggerWrapper) Errorf(format string, args ...interface{}) {
	logger.entry.Errorf(format, args...)
}

func (logger *LoggerWrapper) Fatalf(format string, args ...interface{}) {
	logger.entry.Fatalf(format, args...)
}

func (logger *LoggerWrapper) DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.entry.Debugf(format, args...)
}

func (logger *LoggerWrapper) InfofCtx(ctx context.Context, format string, args ...interface{}) {
	logger.entry.Infof(format, args...)
}

func (logger *LoggerWrapper) WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.entry.Warnf(format, args...)
}

func (logger *LoggerWrapper) ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.entry.Errorf(format, args...)
}

func (logger *LoggerWrapper) FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.entry.Fatalf(format, args...)
}
