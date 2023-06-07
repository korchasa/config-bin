package http

import "net/http"

// handleReadiness is the handler to check if the service is ready to serve requests.
func (s *Server) handleReadiness() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if s.store.IsReady() {
            w.WriteHeader(http.StatusOK)
            _, _ = w.Write([]byte("readiness:OK"))
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
            _, _ = w.Write([]byte("readiness:FAIL"))
        }
    }
}
