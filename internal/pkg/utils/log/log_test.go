package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuncName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Should return current function name",
			expected: "TestGetFuncName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFuncName()
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestLogHandlerInfo(t *testing.T) {
	tests := []struct {
		name       string
		msg        string
		statusCode int
	}{
		{
			name:       "Should log info with status code",
			msg:        "test info message",
			statusCode: 200,
		},
		{
			name:       "Should log info with error status code",
			msg:        "test error info",
			statusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			LogHandlerInfo(logger, tt.msg, tt.statusCode)

			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.msg)
			assert.Contains(t, logOutput, strconv.Itoa(tt.statusCode))
		})
	}
}

func TestLogHandlerError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
		wantErrMsg string
	}{
		{
			name:       "Should log simple error",
			err:        errors.New("test error"),
			statusCode: 500,
			wantErrMsg: "test error",
		},
		{
			name:       "Should log wrapped error",
			err:        fmt.Errorf("wrapper: %w", errors.New("test wrapped error")),
			statusCode: 400,
			wantErrMsg: "test wrapped error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			LogHandlerError(logger, tt.err, tt.statusCode)

			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.wantErrMsg)
			assert.Contains(t, logOutput, strconv.Itoa(tt.statusCode))
		})
	}
}

func TestGetLoggerFromContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		shouldContain string
	}{
		{
			name:          "Should return logger from context",
			ctx:           context.WithValue(context.Background(), "logger", slog.New(slog.NewJSONHandler(os.Stdout, nil))),
			shouldContain: "",
		},
		{
			name:          "Should create new logger when not in context",
			ctx:           context.Background(),
			shouldContain: "Couldnt get logger from context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldContain == "" {
				logger := GetLoggerFromContext(tt.ctx)
				assert.NotNil(t, logger)
			} else {
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				logger := GetLoggerFromContext(tt.ctx)
				assert.NotNil(t, logger)

				w.Close()
				os.Stdout = old

				var output bytes.Buffer
				io.Copy(&output, r)
				assert.Contains(t, output.String(), tt.shouldContain)
			}
		})
	}
}
