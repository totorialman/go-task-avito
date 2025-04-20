package metricsmw

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/totorialman/go-task-avito/internal/pkg/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func CreateHttpMetricsMiddleware(metr *metrics.HttpMetrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := NewResponseWriter(w)
			next.ServeHTTP(rw, r)
			status := http.StatusOK
			statusCode := rw.statusCode
			if statusCode != http.StatusOK && statusCode != http.StatusCreated && statusCode != http.StatusNoContent {
				metr.IncreaseErrors(r.URL.Path, strconv.Itoa(statusCode))
			}
			metr.IncreaseHits(r.URL.Path, strconv.Itoa(statusCode))
			metr.ObserveResponseTime(status, r.URL.Path, time.Since(start).Seconds())
		})
	}
}
