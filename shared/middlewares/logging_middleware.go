package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
)

var loggingLogger = logger.GetLogger()

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		queryParams := r.URL.RawQuery

		// Read and restore request body
		var bodyCopy []byte
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				bodyCopy = bodyBytes
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for next handler
			} else {
				loggingLogger.Debug(fmt.Sprintf("Failed to read request body: %v", err))
			}
		}

		// Wrap response writer
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(lrw, r)

		// Log request and response details
		loggingLogger.Debug("Logger request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("query", queryParams),
			zap.String("body", string(bodyCopy)),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("status", strconv.Itoa(lrw.statusCode)),
			zap.String("time", time.Since(start).String()),
		)
		if lrw.statusCode != 200 {
			loggingLogger.Error("Error response",
				zap.Int("status", lrw.statusCode),
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.String("body", lrw.body.String()),
			)
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(p []byte) (int, error) {
	lrw.body.Write(p) // Capture response body
	return lrw.ResponseWriter.Write(p)
}
