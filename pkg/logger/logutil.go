package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	LogOption func(*logOption)

	logOption struct {
		level     slog.Level
		addSource bool
		// To disable the log writer for STDOut.
		// In default, DisableStdout is false (It appends log lines to stdout.)
		disableStdout bool

		// LogFilename is the file to write logs to.
		// If the LogFilename is empty the logger handler will skip the FS log writer.
		filename            string
		fileRotateMaxSizeMB int
		fileRotateMaxBackup int
		fileRotateCompress  bool

		// customer's logger writer
		logWriter io.Writer
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

func SetDefaultSlog(opts ...LogOption) {
	slog.SetDefault(NewLogger(opts...))
}

func NewLogger(opts ...LogOption) *slog.Logger {
	logOpts := logOption{
		level:               slog.LevelInfo,
		addSource:           false,
		disableStdout:       false,
		filename:            "",
		fileRotateMaxSizeMB: logFileRotateDefaultSizeMB,
		fileRotateMaxBackup: logFileRotateDefaultBackup,
		fileRotateCompress:  false,
		logWriter:           nil,
	}
	for _, opt := range opts {
		opt(&logOpts)
	}
	logWriter := makeLogWriter(logOpts)
	handler := slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
		Level:     logOpts.level,
		AddSource: logOpts.addSource,
	})
	return slog.New(handler)
}

func (nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func LoggerWithSource(addSource bool) LogOption {
	return func(opt *logOption) {
		opt.addSource = addSource
	}
}

func LoggerWithWriter(logWriter io.Writer) LogOption {
	return func(opt *logOption) {
		opt.logWriter = logWriter
	}
}

func LoggerWithStdout(disableStdout bool) LogOption {
	return func(opt *logOption) {
		opt.disableStdout = disableStdout
	}
}

func LoggerWithLevel(levelStr string) LogOption {
	return func(opt *logOption) {
		switch strings.ToLower(strings.TrimSpace(levelStr)) {
		case LOG_LEVEL_DEBUG:
			opt.level = slog.LevelDebug
		case LOG_LEVEL_WARN:
			opt.level = slog.LevelWarn
		case LOG_LEVEL_ERROR:
			opt.level = slog.LevelError
		case LOG_LEVEL_INFO:
			opt.level = slog.LevelInfo
		default:
			opt.level = slog.LevelInfo
		}
	}
}

func LoggerWithFile(
	filename string,
	rotateMaxSizeMB int,
	rotateMaxBackup int,
	rotateCompress bool,
) LogOption {
	return func(opt *logOption) {
		opt.filename = filename
		opt.fileRotateMaxSizeMB = max(rotateMaxSizeMB, logFileRotateMinSizeMB)
		opt.fileRotateMaxBackup = max(rotateMaxBackup, logFileRotateMinBackup)
		opt.fileRotateCompress = rotateCompress
	}
}

func makeLogWriter(opt logOption) io.Writer {
	logFilename := strings.TrimSpace(opt.filename)
	if logFilename == "" && opt.disableStdout && opt.logWriter != nil {
		return &nullWriter{}
	}
	writers := make([]io.Writer, 0, 3)
	if !opt.disableStdout {
		writers = append(writers, os.Stdout)
	}
	if logFilename != "" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    opt.fileRotateMaxSizeMB,
			MaxBackups: opt.fileRotateMaxBackup,
			MaxAge:     0,
			Compress:   opt.fileRotateCompress,
		})
	}
	if opt.logWriter != nil {
		writers = append(writers, opt.logWriter)
	}
	return io.MultiWriter(writers...)
}
