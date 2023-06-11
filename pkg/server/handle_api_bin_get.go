package server

import (
	"configBin/pkg/server/utils"
	"net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleAPIGetBin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bid, err := utils.ExtractBinIDFromPathVar(r)
		if err != nil {
			s.resp.JSONError(r, w, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		pass := utils.PasswordFromHeader(r)
		if pass == "" {
			s.resp.JSONError(r, w, http.StatusBadRequest, "empty_password", err)
			return
		}

		bin, err := s.store.GetBin(*bid, pass)
		if err != nil {
			s.resp.JSONError(r, w, http.StatusBadRequest, "cant_get_bin", err)
			return
		}
		if bin == nil {
			s.resp.JSONError(r, w, http.StatusNotFound, "cant_get_bin", err)
			return
		}

		s.resp.JSONSuccess(w, http.StatusOK, bin)
	}
}
