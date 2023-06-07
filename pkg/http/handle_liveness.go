package http

import "net/http"

// handleLiveness is the handler to check if the service is alive.
func (s *Server) handleLiveness() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("liveness:OK"))
    }
}
