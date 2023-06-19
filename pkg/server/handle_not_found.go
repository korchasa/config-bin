package server

import (
	"errors"
	"net/http"
)

var ErrPageNotFound = errors.New("page not found")

// handleNotFound is the handler to check if the service is alive.
func (s *Server) handleNotFound() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		s.resp.HTMLError(req, resp, http.StatusNotFound, "not_found", ErrPageNotFound)
	}
}
