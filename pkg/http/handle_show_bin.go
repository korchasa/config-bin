package http

import (
    "net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleShowBin() http.HandlerFunc {
    tpl := s.tplProvider.MustGet("bin.gohtml")
    return func(w http.ResponseWriter, r *http.Request) {
        err := tpl.Execute(w, struct {
            Title   string
            Content string
        }{
            Title:   "ConfigBin",
            Content: "Hello, world!",
        })
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}
