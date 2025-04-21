package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendError(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		statusCode   int
		expectedBody map[string]string
	}{
		{
			name:       "Should send 400 error with message",
			message:    "Bad request",
			statusCode: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Bad request",
			},
		},
		{
			name:       "Should send 500 error with message",
			message:    "Internal server error",
			statusCode: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "Internal server error",
			},
		},
		{
			name:       "Should handle empty message",
			message:    "",
			statusCode: http.StatusNotFound,
			expectedBody: map[string]string{
				"error": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			SendError(rr, tt.message, tt.statusCode)

			assert.Equal(t, tt.statusCode, rr.Code)

			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			var responseBody map[string]string
			err := json.NewDecoder(rr.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestSendError_InvalidWriter(t *testing.T) {
	w := &failingWriter{}
	SendError(w, "test error", http.StatusBadRequest)
	assert.True(t, true, "Function should not panic on write error")
}

type failingWriter struct{}

func (f *failingWriter) Header() http.Header {
	return http.Header{}
}

func (f *failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("forced write error")
}

func (f *failingWriter) WriteHeader(statusCode int) {}
