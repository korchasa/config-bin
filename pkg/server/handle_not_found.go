package server

import (
    "fmt"
    "net/http"
)

// handleNotFound is the handler to check if the service is alive.
func (s *Server) handleNotFound() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        s.resp.HTMLError(r, w, http.StatusNotFound, "not_found", fmt.Errorf("page not found"))
    }
}
