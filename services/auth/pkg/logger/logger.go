package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message interface{}, args ...interface{})
	Warn(message interface{}, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zap.SugaredLogger
}

// check
var _ Interface = (*Logger)(nil)

/*func New(level string) *Logger {
	var l zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(l),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			TimeKey:    "ts",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("logger - init: %v", err))
	}

	return &Logger{
		logger: logger.Sugar(),
	}
}*/

func NewWithFile(level string, logFilePath string) Interface {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		l = zapcore.DebugLevel
	case "info":
		l = zapcore.InfoLevel
	case "warn":
		l = zapcore.WarnLevel
	case "error":
		l = zapcore.ErrorLevel
	default:
		l = zapcore.InfoLevel
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(l),
		Encoding:         "json",
		OutputPaths:      []string{"stdout", logFilePath},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "ts",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeName:    zapcore.FullNameEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("logger - init: %v", err))
	}

	return &Logger{
		logger: logger.Sugar(),
	}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.logger.Debugw(fmt.Sprintf("%v", message), args...)
}

func (l *Logger) Info(message interface{}, args ...interface{}) {
	l.logger.Infow(fmt.Sprintf("%v", message), args...)
}

func (l *Logger) Warn(message interface{}, args ...interface{}) {
	l.logger.Warnw(fmt.Sprintf("%v", message), args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.logger.Errorw(fmt.Sprintf("%v", message), args...)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.logger.Fatalw(fmt.Sprintf("%v", message), args...)
}

/*func getServiceName() string {
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}

	return "unknown-service"
}

func getServiceVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}

	return "unknown-version"
}*/
