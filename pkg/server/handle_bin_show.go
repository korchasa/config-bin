package server

import (
	"configBin/pkg"
	"configBin/pkg/server/utils"
	"net/http"
)

// handleBinShow handles bin show request.
func (s *Server) handleBinShow() http.HandlerFunc {
	binTpl := s.tplProvider.MustGet("bin.gohtml")
	authTpl := s.tplProvider.MustGet("auth.gohtml")
	return func(resp http.ResponseWriter, req *http.Request) {
		bid, err := utils.ExtractBinIDFromPathVar(req)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		pass := utils.ReadPasswordCookie(req, *bid)
		if pass == "" {
			bin := pkg.Bin{ID: *bid}
			err = authTpl.Execute(resp, bin)
			if err != nil {
				s.resp.HTMLError(req, resp, http.StatusInternalServerError, "cant_render_template", err)
				return
			}
			return
		}

		bin, err := s.store.GetBin(*bid, pass)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_get_bin", err)
			return
		}
		if bin == nil {
			s.resp.HTMLError(req, resp, http.StatusNotFound, "cant_get_bin", err)
			return
		}

		err = binTpl.Execute(resp, bin)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusInternalServerError, "cant_render_template", err)
			return
		}
		resp.WriteHeader(http.StatusOK)
	}
}
