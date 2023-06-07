package http

import (
    log "github.com/sirupsen/logrus"
    "net/http"
    "time"
)

func (s *Server) logRequest(h http.HandlerFunc) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        start := time.Now()

        uri := r.RequestURI
        method := r.Method
        lrw := NewLoggingResponseWriter(rw)
        s.metrics.IncRequestsCount()
        h(lrw, r)
        duration := time.Since(start)

        if lrw.statusCode >= 200 && lrw.statusCode <= 400 {
            s.metrics.IncSuccessfulRequestsCount(lrw.statusCode)
            s.metrics.ObserveSuccessfulRequestDuration(duration)
        } else {
            s.metrics.IncFailedRequestsCount(lrw.statusCode)
        }
        log.WithFields(
            log.Fields{
                "method":   method,
                "duration": duration,
            },
        ).Infof("Request %s", uri)
    }
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}
