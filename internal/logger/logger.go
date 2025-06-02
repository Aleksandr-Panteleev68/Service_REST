package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(env string) *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if env == "local" {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return &Logger{logger: logger}
}

func (l *Logger) Debug(msg string, keysAndValues ...any){
	l.logger.Debug(msg, keysAndValues...)
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...any) {
	l.logger.Error(msg, append(keysAndValues, "error", err)...)
}

func (l *Logger) Fatal(err error, msg string, keysAndValues ...any) error {
	l.logger.Error(msg, append(keysAndValues, "error", err)...)
	return err
}
