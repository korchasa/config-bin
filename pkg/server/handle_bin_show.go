package server

import (
	"configBin/pkg"
	"configBin/pkg/server/utils"
	"net/http"
)

// handleBinShow handles bin show request
func (s *Server) handleBinShow() http.HandlerFunc {
	binTpl := s.tplProvider.MustGet("bin.gohtml")
	authTpl := s.tplProvider.MustGet("auth.gohtml")
	return func(w http.ResponseWriter, r *http.Request) {
		bid, err := utils.ExtractBIDFromPathVar(r)
		if err != nil {
			s.resp.HTMLError(r, w, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		pass := utils.ReadPassCookie(r, *bid)
		if pass == "" {
			bin := pkg.Bin{ID: *bid}
			err = authTpl.Execute(w, bin)
			if err != nil {
				s.resp.HTMLError(r, w, http.StatusInternalServerError, "cant_render_template", err)
				return
			}
			return
		}

		bin, err := s.store.GetBin(*bid, pass)
		if err != nil {
			s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_get_bin", err)
			return
		}
		if bin == nil {
			s.resp.HTMLError(r, w, http.StatusNotFound, "cant_get_bin", err)
			return
		}

		err = binTpl.Execute(w, bin)
		if err != nil {
			s.resp.HTMLError(r, w, http.StatusInternalServerError, "cant_render_template", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
