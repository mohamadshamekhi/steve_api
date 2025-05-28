package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilePath        = "logs/logs.json"
	defaultUseLocalTime    = false
	defaultFileMaxSizeInMB = 10
	defaultFileAgeInDays   = 30
)

type Config struct {
	FilePath         string
	UseLocalTime     bool
	FileMaxSizeInMB  int
	FileMaxAgeInDays int
}

var l *slog.Logger

// init is default logger and Singleton that lets you ensure that a logger has only one instance, while providing a global access point to this instance.
func init() {
	fileWriter := &lumberjack.Logger{
		Filename:  defaultFilePath,
		LocalTime: defaultUseLocalTime,
		MaxSize:   defaultFileMaxSizeInMB,
		MaxAge:    defaultFileAgeInDays,
	}
	l = slog.New(
		slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), &slog.HandlerOptions{}),
	)
}

func L() *slog.Logger {
	return l
}

// New is constructor logger with special settings.
func New(cfg Config, opt *slog.HandlerOptions) *slog.Logger {
	fileWriter := &lumberjack.Logger{
		Filename:  cfg.FilePath,
		LocalTime: cfg.UseLocalTime,
		MaxSize:   cfg.FileMaxSizeInMB,
		MaxAge:    cfg.FileMaxAgeInDays,
	}

	logger := slog.New(
		slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), opt),
	)

	return logger
}
