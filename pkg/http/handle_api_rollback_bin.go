package http

import (
    "net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleAPIRollbackBin() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        //_, err := w.Write(s.service.CreateBin(r.Context()))
        //if err != nil {
        //    w.WriteHeader(http.StatusInternalServerError)
        //    return
        //}
        w.WriteHeader(http.StatusOK)
    }
}
