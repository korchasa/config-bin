package server

import (
    "configBin/pkg/server/utils"
    "fmt"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
    "net/http"
)

// handleBinCreate is the handler to create a bin.
func (s *Server) handleBinCreate() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_parse_form", err)
            return
        }

        id := r.Form.Get("uuid")
        if id == "" {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "empty_uuid", fmt.Errorf("empty uuid"))
            return
        }
        bid, err := uuid.Parse(id)
        if err != nil {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "invalid_bin_id", fmt.Errorf("invalid bin id"))
            return
        }

        pass := r.Form.Get("password")
        if pass == "" {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "empty_password", fmt.Errorf("empty password"))
            return
        }

        unencryptedData := r.Form.Get("content")
        if unencryptedData == "" {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "empty_content", fmt.Errorf("empty content"))
        }

        err = s.store.CreateBin(bid, pass, unencryptedData)
        if err != nil {
            s.resp.HTMLError(r, w, http.StatusBadRequest, "cant_create_bin", err)
            return
        }

        utils.WritePassCookie(w, bid, pass)

        log.Infof("bin created: %s", bid.String())

        http.Redirect(w, r, fmt.Sprintf("/%s", bid.String()), http.StatusFound)
    }
}
