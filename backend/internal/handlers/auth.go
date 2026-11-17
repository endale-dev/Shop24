package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/domain/models"
	"shop24/internal/httpx"
	"shop24/internal/services"
	"shop24/internal/token"
)

type AuthHandler struct {
	auth         *services.AuthService
	tokens       *token.Issuer
	log          *slog.Logger
	secure       bool
	domain       string
	publicURL    string
	googleID     string
	googleSecret string
}

func NewAuthHandler(auth *services.AuthService, tokens *token.Issuer, log *slog.Logger, secure bool, domain, publicURL, googleID, googleSecret string) *AuthHandler {
	return &AuthHandler{
		auth: auth, tokens: tokens, log: log, secure: secure, domain: domain,
		publicURL: publicURL, googleID: googleID, googleSecret: googleSecret,
	}
}

func (h *AuthHandler) Providers(w http.ResponseWriter, r *http.Request) {
	httpx.Data(w, http.StatusOK, map[string]bool{"google": h.googleID != "" && h.googleSecret != ""})
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
	h.writeUser(w, http.StatusCreated, user)
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
	h.writeUser(w, http.StatusOK, user)
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
	h.writeUser(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	user, err := h.auth.UpdateName(r.Context(), id, body.Name)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	h.writeUser(w, http.StatusOK, user)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Current  string `json:"current"`
		Password string `json:"password"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	user, err := h.auth.ChangePassword(r.Context(), id, body.Current, body.Password)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	h.writeUser(w, http.StatusOK, user)
}

func (h *AuthHandler) Forgot(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	if err := h.auth.RequestReset(r.Context(), body.Email); err != nil {
		if appErr, ok := apperr.As(err); ok && appErr.Status < 500 {
			fail(h.log, w, err)
			return
		}
		h.log.Error("password reset", "error", err)
	}
	httpx.Data(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) Reset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	user, signed, exp, err := h.auth.ResetPassword(r.Context(), body.Token, body.Password)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	h.setSession(w, signed, exp)
	h.writeUser(w, http.StatusOK, user)
}

func (h *AuthHandler) writeUser(w http.ResponseWriter, status int, user models.User) {
	user.HasPassword = user.PasswordHash != ""
	user.GoogleLinked = user.GoogleSub != ""
	httpx.Data(w, status, user)
}

func (h *AuthHandler) setSession(w http.ResponseWriter, value string, exp time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name: token.CookieName, Value: value, Path: "/", Expires: exp,
		MaxAge: int(time.Until(exp).Seconds()), HttpOnly: true, Secure: h.secure,
		SameSite: http.SameSiteLaxMode, Domain: h.domain,
	})
}
