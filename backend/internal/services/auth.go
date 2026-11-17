package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"shop24/internal/apperr"
	"shop24/internal/domain/input"
	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
	"shop24/internal/token"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type AuthService struct {
	users     interfaces.UserRepository
	tokens    *token.Issuer
	mail      Mailer
	publicURL string
}

func NewAuthService(users interfaces.UserRepository, tokens *token.Issuer, mail Mailer, publicURL string) *AuthService {
	return &AuthService{
		users: users, tokens: tokens, mail: mail,
		publicURL: strings.TrimRight(strings.TrimSpace(publicURL), "/"),
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (models.User, string, time.Time, error) {
	name, err := input.Text(name, "Name", 2, 80)
	if err != nil {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusBadRequest, "invalid_name", err.Error())
	}
	email, err = input.Email(email)
	if err != nil {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusBadRequest, "invalid_email", err.Error())
	}
	if err := validPassword(password); err != nil {
		return models.User{}, "", time.Time{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	now := time.Now().UTC()
	user, err := s.users.Create(ctx, models.User{
		Name: name, Email: email, PasswordHash: string(hash), CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.User{}, "", time.Time{}, appErr
		}
		return models.User{}, "", time.Time{}, err
	}
	return s.issue(user, now)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (models.User, string, time.Time, error) {
	email, err := input.Email(email)
	if err != nil {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusUnauthorized, "invalid_login", "Email or password is incorrect")
	}
	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusUnauthorized, "invalid_login", "Email or password is incorrect")
	}
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	if user.PasswordHash == "" || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusUnauthorized, "invalid_login", "Email or password is incorrect")
	}
	return s.issue(user, time.Now().UTC())
}

func (s *AuthService) Current(ctx context.Context, userID string) (models.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.User{}, apperr.ErrUnauthorized
	}
	return user, err
}

func (s *AuthService) LoginGoogle(ctx context.Context, sub, email, name string) (models.User, string, time.Time, error) {
	email, err := input.Email(email)
	if err != nil || strings.TrimSpace(sub) == "" {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusBadRequest, "google", "Google did not return a usable account")
	}
	sub = strings.TrimSpace(sub)
	existing, err := s.users.FindByGoogleSub(ctx, sub)
	if err == nil {
		return s.issue(existing, time.Now().UTC())
	}
	if !errors.Is(err, interfaces.ErrNotFound) {
		return models.User{}, "", time.Time{}, err
	}
	byEmail, err := s.users.FindByEmail(ctx, email)
	if err == nil {
		if byEmail.GoogleSub != "" && byEmail.GoogleSub != sub {
			return models.User{}, "", time.Time{}, apperr.New(http.StatusConflict, "google_taken", "This email is already linked to a different Google account")
		}
		byEmail.GoogleSub = sub
		byEmail.UpdatedAt = time.Now().UTC()
		if err := s.users.Update(ctx, byEmail); err != nil {
			if appErr, ok := conflict(err); ok {
				return models.User{}, "", time.Time{}, appErr
			}
			return models.User{}, "", time.Time{}, err
		}
		return s.issue(byEmail, byEmail.UpdatedAt)
	}
	if !errors.Is(err, interfaces.ErrNotFound) {
		return models.User{}, "", time.Time{}, err
	}
	display, err := sellerName(name, email)
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	now := time.Now().UTC()
	user, err := s.users.Create(ctx, models.User{
		Name: display, Email: email, GoogleSub: sub, CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.User{}, "", time.Time{}, appErr
		}
		return models.User{}, "", time.Time{}, err
	}
	return s.issue(user, now)
}

func (s *AuthService) LinkGoogle(ctx context.Context, userID, sub, email string) (models.User, error) {
	email, err := input.Email(email)
	if err != nil || strings.TrimSpace(sub) == "" {
		return models.User{}, apperr.New(http.StatusBadRequest, "google", "Google did not return a usable account")
	}
	user, err := s.Current(ctx, userID)
	if err != nil {
		return models.User{}, err
	}
	if user.Email != email {
		return models.User{}, apperr.New(http.StatusBadRequest, "google_email", "Google must use the same email as this account")
	}
	sub = strings.TrimSpace(sub)
	other, err := s.users.FindByGoogleSub(ctx, sub)
	if err == nil && other.ID != user.ID {
		return models.User{}, apperr.New(http.StatusConflict, "google_taken", "That Google account is already used")
	}
	if err != nil && !errors.Is(err, interfaces.ErrNotFound) {
		return models.User{}, err
	}
	user.GoogleSub = sub
	user.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, user); err != nil {
		if appErr, ok := conflict(err); ok {
			return models.User{}, appErr
		}
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) UpdateName(ctx context.Context, userID, name string) (models.User, error) {
	name, err := input.Text(name, "Name", 2, 80)
	if err != nil {
		return models.User{}, apperr.New(http.StatusBadRequest, "invalid_name", err.Error())
	}
	user, err := s.Current(ctx, userID)
	if err != nil {
		return models.User{}, err
	}
	user.Name = name
	user.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, current, next string) (models.User, error) {
	if err := validPassword(next); err != nil {
		return models.User{}, err
	}
	user, err := s.Current(ctx, userID)
	if err != nil {
		return models.User{}, err
	}
	if user.PasswordHash != "" && bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(current)) != nil {
		return models.User{}, apperr.New(http.StatusBadRequest, "invalid_password", "Current password is incorrect")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.PasswordHash = string(hash)
	user.ResetHash = ""
	user.ResetExpires = time.Time{}
	user.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) RequestReset(ctx context.Context, email string) error {
	email, err := input.Email(email)
	if err != nil {
		return apperr.New(http.StatusBadRequest, "invalid_email", "Enter a valid email")
	}
	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, interfaces.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return err
	}
	tokenValue := hex.EncodeToString(raw)
	user.ResetHash = hashToken(tokenValue)
	user.ResetExpires = time.Now().UTC().Add(30 * time.Minute)
	user.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, user); err != nil {
		return err
	}
	body := "Reset the password for your Sable account.\n\n" + s.publicURL + "/reset?token=" + tokenValue + "\n\nThis link stops working in 30 minutes. If you did not ask for it, ignore this email.\n"
	if s.mail == nil {
		return errors.New("mail is not configured")
	}
	return s.mail.Send(ctx, user.Email, "Reset your Sable password", body)
}

func (s *AuthService) ResetPassword(ctx context.Context, tokenValue, password string) (models.User, string, time.Time, error) {
	if err := validPassword(password); err != nil {
		return models.User{}, "", time.Time{}, err
	}
	user, err := s.users.FindByResetHash(ctx, hashToken(strings.TrimSpace(tokenValue)))
	if errors.Is(err, interfaces.ErrNotFound) || tokenValue == "" {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusBadRequest, "invalid_reset", "This reset link is not valid")
	}
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	if user.ResetExpires.Before(time.Now().UTC()) {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusBadRequest, "invalid_reset", "This reset link is not valid")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	now := time.Now().UTC()
	user.PasswordHash = string(hash)
	user.ResetHash = ""
	user.ResetExpires = time.Time{}
	user.UpdatedAt = now
	if err := s.users.Update(ctx, user); err != nil {
		return models.User{}, "", time.Time{}, err
	}
	return s.issue(user, now)
}

func (s *AuthService) issue(user models.User, now time.Time) (models.User, string, time.Time, error) {
	signed, exp, err := s.tokens.Issue(user.ID, now)
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	return user, signed, exp, nil
}

func validPassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return apperr.New(http.StatusBadRequest, "invalid_password", "Use a password between 8 and 72 characters")
	}
	return nil
}

func hashToken(tokenValue string) string {
	sum := sha256.Sum256([]byte(tokenValue))
	return hex.EncodeToString(sum[:])
}

func sellerName(name, email string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = strings.Split(email, "@")[0]
	}
	cleaned, err := input.Text(name, "Name", 2, 80)
	if err != nil {
		return input.Text("Seller", "Name", 2, 80)
	}
	return cleaned, nil
}
