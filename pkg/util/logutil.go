package util

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	LogOption struct {
		Level     string
		AddSource bool
		// To disable the log writer for STDOut.
		// In default, DisableStdout is false (It appends log lines to stdout.)
		DisableStdout bool

		// LogFilename is the file to write logs to.
		// If the LogFilename is empty the logger handler will skip the FS log writer.
		LogFilename            string
		LogFileRotateMaxSizeMB int
		LogFileRotateMaxBackup int
		LogFileRtateCompress   bool
	}

	nullWriter struct{}
)

const (
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"

	logFileRotateMinSizeMB     = 5
	logFileRotateDefaultSizeMB = 100
	logFileRotateMinBackup     = 1
	logFileRotateDefaultBackup = 5
)

func NewLogger(opts LogOption) *slog.Logger {
	return slog.New(makeJsonLogHandler(opts))
}

func (nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case LOG_LEVEL_DEBUG:
		return slog.LevelDebug
	case LOG_LEVEL_WARN:
		return slog.LevelWarn
	case LOG_LEVEL_ERROR:
		return slog.LevelError
	case LOG_LEVEL_INFO:
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}

func parseLogFileRotateMaxBackup(opts LogOption) int {
	if opts.LogFileRotateMaxBackup == 0 {
		return logFileRotateDefaultBackup
	} else if opts.LogFileRotateMaxBackup < logFileRotateMinBackup {
		return logFileRotateMinBackup
	} else {
		return opts.LogFileRotateMaxBackup
	}
}

func parseLogFileRotateMaxSize(opts LogOption) int {
	if opts.LogFileRotateMaxSizeMB == 0 {
		return logFileRotateDefaultSizeMB
	} else if opts.LogFileRotateMaxSizeMB < logFileRotateMinSizeMB {
		return logFileRotateMinSizeMB
	} else {
		return opts.LogFileRotateMaxSizeMB
	}
}

func makeLogWriter(opts LogOption) io.Writer {
	logFilename := strings.TrimSpace(opts.LogFilename)
	if logFilename == "" && opts.DisableStdout {
		return &nullWriter{}
	}

	writers := make([]io.Writer, 0, 2)
	if !opts.DisableStdout {
		writers = append(writers, os.Stdout)
	}
	if logFilename != "" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    parseLogFileRotateMaxSize(opts),
			MaxBackups: parseLogFileRotateMaxBackup(opts),
			MaxAge:     0,
			Compress:   opts.LogFileRtateCompress,
		})
	}
	return io.MultiWriter(writers...)
}

func makeJsonLogHandler(opts LogOption) slog.Handler {
	handlerOpts := &slog.HandlerOptions{
		Level:     parseLogLevel(opts.Level),
		AddSource: opts.AddSource,
	}
	logWriter := makeLogWriter(opts)
	return slog.NewJSONHandler(logWriter, handlerOpts)
}
