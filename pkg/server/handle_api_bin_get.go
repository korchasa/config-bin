package server

import (
	"configBin/pkg/server/utils"
	"net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleAPIGetBin() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		bid, err := utils.ExtractBinIDFromPathVar(req)
		if err != nil {
			s.resp.JSONError(req, resp, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		pass := utils.PasswordFromHeader(req)
		if pass == "" {
			s.resp.JSONError(req, resp, http.StatusBadRequest, "empty_password", err)
			return
		}

		bin, err := s.store.GetBin(*bid, pass)
		if err != nil {
			s.resp.JSONError(req, resp, http.StatusBadRequest, "cant_get_bin", err)
			return
		}
		if bin == nil {
			s.resp.JSONError(req, resp, http.StatusNotFound, "cant_get_bin", err)
			return
		}

		s.resp.JSONSuccess(resp, http.StatusOK, bin)
	}
}
