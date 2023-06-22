package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *Server) logRequest(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		uri := req.RequestURI
		method := req.Method
		lrw := NewLoggingResponseWriter(resp)
		s.metrics.IncRequestsCount()
		handlerFunc(lrw, req)
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

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
