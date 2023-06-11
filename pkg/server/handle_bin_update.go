package server

import (
	"configBin/pkg/server/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// handleBinUpdate is the handler to update a bin.
func (s *Server) handleBinUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bid, err := utils.ExtractBinIDFromPathVar(r)
		if err != nil {
			s.resp.HTMLError(r, w, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		pass := utils.ReadPasswordCookie(r, *bid)

		err = r.ParseForm()
		if err != nil {
			s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_parse_form", err)
			return
		}

		unencryptedData := r.Form.Get("content")
		if unencryptedData == "" {
			s.resp.HTMLError(r, w, http.StatusBadRequest, "no_content", fmt.Errorf("content is empty"))
			return
		}

		err = s.store.UpdateBin(*bid, pass, unencryptedData)
		if err != nil {
			s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_update_bin", err)
			return
		}

		log.Infof("bin updated: %s", bid.String())

		http.Redirect(w, r, fmt.Sprintf("/%s", bid.String()), http.StatusFound)
	}
}
