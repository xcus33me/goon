package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message interface{}, args ...interface{})
	Warn(message interface{}, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *slog.Logger
}

// check
var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l slog.Level

	switch strings.ToLower(level) {
	case "error":
		l = slog.LevelError
	case "warn":
		l = slog.LevelWarn
	case "info":
		l = slog.LevelInfo
	case "debug":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     l,
		AddSource: true,
	})

	logger := slog.New(handler)

	return &Logger{logger: logger}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.log(slog.LevelDebug, message, args...)
}

func (l *Logger) Info(message interface{}, args ...interface{}) {
	l.log(slog.LevelInfo, message, args...)
}

func (l *Logger) Warn(message interface{}, args ...interface{}) {
	l.log(slog.LevelWarn, message, args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.log(slog.LevelError, message, args...)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.log(slog.LevelError, message, args...)
	os.Exit(1)
}

func (l *Logger) log(level slog.Level, message interface{}, args ...interface{}) {
	if !l.logger.Enabled(nil, level) {
		return
	}

	switch msg := message.(type) {
	case string:
		if len(args) == 0 {
			l.logger.Log(nil, level, msg)
		} else {
			l.logger.Log(nil, level, fmt.Sprintf(msg, args...))
		}
	case error:
		l.logger.Log(nil, level, msg.Error(), slog.Any("error", msg))
	default:
		l.logger.Log(nil, level, fmt.Sprintf("unknown message type: %v", msg), slog.Any("value", msg))
	}
}
