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
		Level         string
		AddSource     bool
		DisableStdout bool

		LogFilename            string
		LogFileRotateMaxSizeMB int
		LogFileRotateMaxBackup int
		LogFileRtateCompress   bool
	}

	logNullWriter struct{}
)

const (
	LOG_FILE_ROTATE_MIN_SIZE_MB     = 5
	LOG_FILE_ROTATE_DEFAULT_SIZE_MB = 300
	LOG_FILE_ROTATE_MIN_BACKUP      = 1
	LOG_FILE_ROTATE_DEFAULT_BACKUP  = 5
)

func (logNullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func parseLogFileRotateMaxBackup(opts LogOption) int {
	if opts.LogFileRotateMaxBackup == 0 {
		return LOG_FILE_ROTATE_DEFAULT_BACKUP
	} else if opts.LogFileRotateMaxBackup < LOG_FILE_ROTATE_MIN_BACKUP {
		return LOG_FILE_ROTATE_MIN_BACKUP
	} else {
		return opts.LogFileRotateMaxBackup
	}
}

func parseLogFileRotateMaxSize(opts LogOption) int {
	if opts.LogFileRotateMaxSizeMB == 0 {
		return LOG_FILE_ROTATE_DEFAULT_SIZE_MB
	} else if opts.LogFileRotateMaxSizeMB < LOG_FILE_ROTATE_MIN_SIZE_MB {
		return LOG_FILE_ROTATE_MIN_SIZE_MB
	} else {
		return opts.LogFileRotateMaxSizeMB
	}
}

func makeLogWriter(opts LogOption) io.Writer {
	logFilename := strings.TrimSpace(opts.LogFilename)
	if logFilename == "" && opts.DisableStdout {
		return &logNullWriter{}
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

func NewLog(opts LogOption) *slog.Logger {
	return slog.New(makeJsonLogHandler(opts))
}

func InitDefaultLog(opts LogOption) {
	slog.SetDefault(NewLog(opts))
}
