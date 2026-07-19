package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	*slog.Logger
	file     *os.File
	filePath string
}

func NewLogger(config Config) (*Logger, error) {

	level, err := ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logger level name: %w", err)
	}

	var handlers []slog.Handler

	var (
		logFile     *os.File
		logFilePath string
	)

	if config.Folder != "" {
		if err := os.MkdirAll(config.Folder, 0755); err != nil {
			return nil, fmt.Errorf("failed to make logger folder: %w", err)
		}
		logFile, logFilePath, err = newLogFile(config.Folder)
		if err != nil {
			return nil, fmt.Errorf("failed to open logfile: %w", err)
		}

		handlers = append(handlers, slog.NewJSONHandler(
			logFile,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     level,
			},
		))
	}

	handlers = append(handlers, slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: level,
		},
	))

	logger := slog.New(newMultiHandler(handlers...))

	return &Logger{
		Logger:   logger,
		file:     logFile,
		filePath: logFilePath,
	}, nil
}

func newLogFile(folder string) (*os.File, string, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(folder, fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, "", err
	}

	return logFile, logFilePath, nil
}

func ParseLevel(s string) (slog.Level, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level name %q (want debug|info|warn|error)", s)
	}
}

func (l *Logger) Close() error {
	if l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *Logger) FilePath() string {
	return l.filePath
}
