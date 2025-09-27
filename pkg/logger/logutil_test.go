package logger

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogLevel(t *testing.T) {
	assert := require.New(t)

	l := NewLogger(LoggerWithLevel("info"))
	assert.False(l.Enabled(t.Context(), slog.LevelDebug), "Level debug is disabled")
	assert.True(l.Enabled(t.Context(), slog.LevelInfo), "Level info is enable")
}

func TestLogWriter(t *testing.T) {
	assert := require.New(t)

	var logBuffer bytes.Buffer
	l := NewLogger(
		LoggerWithLevel("info"),
		LoggerWithWriter(&logBuffer),
	)

	l.Info("test log", slog.String("msg", "value"))
	logData := logBuffer.String()

	t.Logf("LogData: %s\n", logData)
	assert.Contains(logData, "INFO")
	assert.Contains(logData, "test log")
	assert.Contains(logData, "value")
}
