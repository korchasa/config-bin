package http

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "net/http"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleBinAuth() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        urlVars := mux.Vars(r)
        bid, exists := urlVars["bid"]
        if !exists {
            s.handleError(r, w, http.StatusBadRequest, "no_bin_id", fmt.Errorf("no bin id"))
            return
        }
        id, err := uuid.Parse(bid)
        if err != nil {
            s.handleError(r, w, http.StatusBadRequest, "invalid_bin_id", fmt.Errorf("invalid bin id"))
            return
        }

        err = r.ParseForm()
        if err != nil {
            s.handleError(r, w, http.StatusBadRequest, "cant_parse_form", err)
            return
        }

        pass := r.Form.Get("password")
        if pass == "" {
            s.handleError(r, w, http.StatusBadRequest, "empty_password", fmt.Errorf("empty password"))
            return
        }

        unencryptedData := r.Form.Get("text")

        err = s.store.CreateBin(id, pass, unencryptedData)
        if err != nil {
            s.handleError(r, w, http.StatusBadRequest, "cant_create_bin", nil)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     id.String(),
            Value:    pass,
            Path:     "/",
            MaxAge:   3600,
            HttpOnly: true,
            SameSite: http.SameSiteLaxMode,
        })

        log.Infof("bin created: %s", id.String())

        http.Redirect(w, r, fmt.Sprintf("/%s", id.String()), http.StatusFound)
    }
}
