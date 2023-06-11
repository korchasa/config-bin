package server

import (
    "net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleAPIGetBin() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }
}
