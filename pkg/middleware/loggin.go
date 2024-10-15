package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggingMiddleware struct {
	lg *zap.Logger
}

func NewLogging(lg *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		lg: lg,
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (lm *LoggingMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		lm.lg.Log(
			zapcore.InfoLevel,
			"request response",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
			zap.Duration("Time", time.Since(start)),
		)
	})
}
