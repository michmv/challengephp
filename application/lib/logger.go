package lib

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Close()
}

type FileLogger struct {
	file *os.File
	*slog.Logger
}

func NewFileLogger(path string) (FileLogger, Error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return FileLogger{}, Err(err)
	}

	logger := slog.New(slog.NewJSONHandler(file, nil))

	return FileLogger{
		file:   file,
		Logger: logger,
	}, nil
}

func (it FileLogger) Close() {
	_ = it.file.Close()
}
