package utils

import (
	"github.com/google/uuid"
	"net/http"
)

func PasswordCookie(bid uuid.UUID, pass string) *http.Cookie {
	return &http.Cookie{
		Name:     bid.String(),
		Value:    pass,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func ReadPasswordCookie(r *http.Request, bid uuid.UUID) string {
	cookie, err := r.Cookie(bid.String())
	if err != nil {
		return ""
	}
	return cookie.Value
}

const PasswordHeader = "X-ConfigBin-Password"

func PasswordFromHeader(r *http.Request) string {
	return r.Header.Get(PasswordHeader)
}
