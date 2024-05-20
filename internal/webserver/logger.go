package webserver

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogging(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriterWrapper{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		elapsed := time.Since(start)

		logger.Printf("%s %s | Response: %d | Time: %s", r.Method, r.URL.String(), rw.statusCode, elapsed.String())
	})
}
