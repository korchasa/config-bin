package server

import (
    "configBin/pkg/server/utils"
    "fmt"
    log "github.com/sirupsen/logrus"
    "net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleBinAuth() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        bid, err := utils.ExtractBIDFromPathVar(r)
        if err != nil {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "invalid_bin_id", err)
            return
        }

        err = r.ParseForm()
        if err != nil {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_parse_form", err)
            return
        }

        pass := r.Form.Get("password")
        if pass == "" {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "empty_password", fmt.Errorf("empty password"))
            return
        }

        bin, err := s.store.GetBin(*bid, pass)
        if err != nil {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_get_bin", err)
            return
        }
        if bin == nil {
            s.resp.HTMLError(r, w, http.StatusNotFound, "no_bin_by_id", fmt.Errorf("no bin found by id"))
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     bid.String(),
            Value:    pass,
            Path:     "/",
            MaxAge:   3600,
            HttpOnly: true,
            SameSite: http.SameSiteLaxMode,
        })

        log.Infof("bin authed: %s", bid.String())

        http.Redirect(w, r, fmt.Sprintf("/%s", bid.String()), http.StatusFound)
    }
}
