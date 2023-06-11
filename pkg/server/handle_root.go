package server

import (
    "github.com/google/uuid"
    "math/rand"
    "net/http"
    "strings"
    "time"
)

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleRoot() http.HandlerFunc {
    tpl := s.tplProvider.MustGet("root.gohtml")
    return func(w http.ResponseWriter, r *http.Request) {
        err := tpl.Execute(w, struct {
            ID       string
            Password string
        }{
            ID:       generateUUID(),
            Password: generatePassword(8),
        })
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}

func generateUUID() string {
    return uuid.New().String()
}

// generatePassword generates a password of the given length
func generatePassword(length int) string {
    upperCaseLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    lowerCaseLetters := strings.ToLower(upperCaseLetters)
    digits := "0123456789"
    specials := "!@#$%&*()-_+=~"
    all := upperCaseLetters + lowerCaseLetters + digits + specials

    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    var password strings.Builder
    for i := 0; i < length; i++ {
        randomIndex := r.Intn(len(all))
        password.WriteByte(all[randomIndex])
    }
    return password.String()
}
