package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))
	return &Logger{logger: logger}
}

func (l *Logger) Info(msg string, keysAndValues ...any){
	l.logger.Info(msg, keysAndValues...)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...any){
	l.logger.Error(msg, append(keysAndValues, "error", err)...)
}

func (l *Logger) Fatal(err error, msg string, keysAndValues ...any){
	l.logger.Error(msg, append(keysAndValues, "error", err)...)
	os.Exit(1)
}