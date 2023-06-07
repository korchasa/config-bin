package http

import "net/http"

// handleReadiness is the handler to check if the service is ready to serve requests.
func (s *Server) handleReadiness() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if s.store.IsReady() {
            s.sendOKResponse(w, http.StatusOK, "readiness:OK")
        } else {
            s.sendOKResponse(w, http.StatusServiceUnavailable, "readiness:FAIL")
        }
    }
}
