package logger

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	l := NewLogger(LoggerWithLevel("info"))
	assert.False(t, l.Enabled(t.Context(), slog.LevelDebug), "Level debug is disabled")
	assert.True(t, l.Enabled(t.Context(), slog.LevelInfo), "Level info is enable")
}

func TestLogWriter(t *testing.T) {
	var logBuffer bytes.Buffer
	l := NewLogger(
		LoggerWithLevel("info"),
		LoggerWithWriter(&logBuffer),
	)

	l.Info("test log", slog.String("msg", "value"))
	logData := logBuffer.String()

	t.Logf("LogData: %s\n", logData)
	assert.Contains(t, logData, "INFO")
	assert.Contains(t, logData, "test log")
	assert.Contains(t, logData, "value")
}
