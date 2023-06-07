package http

import (
    "fmt"
    "net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleAPICreateBin() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        //_, err := w.Write(s.service.CreateBin(r.Context()))
        //if err != nil {
        //    w.WriteHeader(http.StatusInternalServerError)
        //    return
        //}
        if r.Method == "POST" {
            s.handleError(r, w, http.StatusMethodNotAllowed, "method_not_allowed", fmt.Errorf("method not allowed"))
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}
