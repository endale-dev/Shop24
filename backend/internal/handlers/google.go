package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/token"
)

const oauthCookie = "sable_oauth"

func (h *AuthHandler) GoogleStart(w http.ResponseWriter, r *http.Request) {
	if h.googleID == "" || h.googleSecret == "" {
		h.redirect(w, r, "/login", "google")
		return
	}
	mode := "login"
	if r.URL.Query().Get("mode") == "link" {
		mode = "link"
	}
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		h.redirect(w, r, h.failPath(mode), "google")
		return
	}
	state := hex.EncodeToString(raw)
	http.SetCookie(w, &http.Cookie{
		Name: oauthCookie, Value: state + "." + mode, Path: "/", MaxAge: 600,
		HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteLaxMode, Domain: h.domain,
	})
	query := url.Values{}
	query.Set("client_id", h.googleID)
	query.Set("redirect_uri", h.publicURL+"/api/auth/google/callback")
	query.Set("response_type", "code")
	query.Set("scope", "openid email profile")
	query.Set("state", state)
	query.Set("prompt", "select_account")
	http.Redirect(w, r, "https://accounts.google.com/o/oauth2/v2/auth?"+query.Encode(), http.StatusFound)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state, mode := readOAuthCookie(r)
	http.SetCookie(w, &http.Cookie{
		Name: oauthCookie, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteLaxMode, Domain: h.domain,
	})
	failPath := h.failPath(mode)
	if state == "" || state != r.URL.Query().Get("state") || r.URL.Query().Get("error") != "" {
		h.redirect(w, r, failPath, "google")
		return
	}
	profile, err := exchangeGoogle(r.Context(), h.googleID, h.googleSecret, h.publicURL+"/api/auth/google/callback", r.URL.Query().Get("code"))
	if err != nil {
		h.log.Error("google sign-in", "error", err)
		code := "google"
		if err.Error() == "unverified" {
			code = "google_unverified"
		}
		h.redirect(w, r, failPath, code)
		return
	}
	if mode == "link" {
		id, err := h.sessionID(r)
		if err != nil {
			h.redirect(w, r, "/login", "google")
			return
		}
		if _, err := h.auth.LinkGoogle(r.Context(), id, profile.Sub, profile.Email); err != nil {
			code := "google"
			if appErr, ok := apperr.As(err); ok && appErr.Code != "" {
				code = appErr.Code
			}
			h.redirect(w, r, "/account", code)
			return
		}
		h.redirect(w, r, "/account", "")
		return
	}
	_, signed, exp, err := h.auth.LoginGoogle(r.Context(), profile.Sub, profile.Email, profile.Name)
	if err != nil {
		code := "google"
		if appErr, ok := apperr.As(err); ok && appErr.Code != "" {
			code = appErr.Code
		}
		h.redirect(w, r, "/login", code)
		return
	}
	h.setSession(w, signed, exp)
	h.redirect(w, r, "/", "")
}

func (h *AuthHandler) failPath(mode string) string {
	if mode == "link" {
		return "/account"
	}
	return "/login"
}

func (h *AuthHandler) redirect(w http.ResponseWriter, r *http.Request, path, code string) {
	target := h.publicURL + path
	if code != "" {
		target += "?error=" + url.QueryEscape(code)
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *AuthHandler) sessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(token.CookieName)
	if err != nil || cookie.Value == "" {
		return "", apperr.ErrUnauthorized
	}
	return h.tokens.Parse(cookie.Value)
}

func readOAuthCookie(r *http.Request) (string, string) {
	cookie, err := r.Cookie(oauthCookie)
	if err != nil {
		return "", "login"
	}
	state, mode, ok := strings.Cut(cookie.Value, ".")
	if !ok || (mode != "login" && mode != "link") || len(state) < 16 {
		return "", "login"
	}
	return state, mode
}

type googleProfile struct {
	Sub   string
	Email string
	Name  string
}

func exchangeGoogle(ctx context.Context, clientID, secret, redirectURI, code string) (googleProfile, error) {
	if code == "" {
		return googleProfile{}, errGoogle
	}
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return googleProfile{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 12 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return googleProfile{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return googleProfile{}, err
	}
	var tokenBody struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenBody); err != nil || tokenBody.AccessToken == "" || res.StatusCode >= 300 {
		return googleProfile{}, errGoogle
	}
	infoReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return googleProfile{}, err
	}
	infoReq.Header.Set("Authorization", "Bearer "+tokenBody.AccessToken)
	infoRes, err := client.Do(infoReq)
	if err != nil {
		return googleProfile{}, err
	}
	defer infoRes.Body.Close()
	infoBody, err := io.ReadAll(io.LimitReader(infoRes.Body, 1<<20))
	if err != nil {
		return googleProfile{}, err
	}
	var raw map[string]any
	if err := json.Unmarshal(infoBody, &raw); err != nil || infoRes.StatusCode >= 300 {
		return googleProfile{}, errGoogle
	}
	if !googleVerified(raw["email_verified"]) {
		return googleProfile{}, errUnverified
	}
	profile := googleProfile{Sub: googleString(raw["sub"]), Email: googleString(raw["email"]), Name: googleString(raw["name"])}
	if profile.Sub == "" || profile.Email == "" {
		return googleProfile{}, errGoogle
	}
	return profile, nil
}

func googleVerified(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return typed == "true"
	default:
		return false
	}
}

func googleString(value any) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
}

var errGoogle = googleError("google")
var errUnverified = googleError("unverified")

type googleError string

func (e googleError) Error() string { return string(e) }
