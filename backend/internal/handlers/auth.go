package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"shop24/internal/httpx"
	"shop24/internal/services"
	"shop24/internal/token"
)

type AuthHandler struct {
	auth   *services.AuthService
	log    *slog.Logger
	secure bool
	domain string
}

func NewAuthHandler(auth *services.AuthService, log *slog.Logger, secure bool, domain string) *AuthHandler {
	return &AuthHandler{auth: auth, log: log, secure: secure, domain: domain}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	user, signed, exp, err := h.auth.Register(r.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	h.setSession(w, signed, exp)
	httpx.Data(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	user, signed, exp, err := h.auth.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	h.setSession(w, signed, exp)
	httpx.Data(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: token.CookieName, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteLaxMode, Domain: h.domain,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	id, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	user, err := h.auth.Current(r.Context(), id)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, user)
}

func (h *AuthHandler) setSession(w http.ResponseWriter, value string, exp time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name: token.CookieName, Value: value, Path: "/", Expires: exp,
		MaxAge: int(time.Until(exp).Seconds()), HttpOnly: true, Secure: h.secure,
		SameSite: http.SameSiteLaxMode, Domain: h.domain,
	})
}
