package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"shop24/internal/apperr"
	"shop24/internal/domain/input"
	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
	"shop24/internal/token"
)

type AuthService struct {
	users  interfaces.UserRepository
	tokens *token.Issuer
}

func NewAuthService(users interfaces.UserRepository, tokens *token.Issuer) *AuthService {
	return &AuthService{users: users, tokens: tokens}
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
	if len(password) < 8 || len(password) > 72 {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusBadRequest, "invalid_password", "Use a password between 8 and 72 characters")
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
	signed, exp, err := s.tokens.Issue(user.ID, now)
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	return user, signed, exp, nil
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
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return models.User{}, "", time.Time{}, apperr.New(http.StatusUnauthorized, "invalid_login", "Email or password is incorrect")
	}
	now := time.Now().UTC()
	signed, exp, err := s.tokens.Issue(user.ID, now)
	if err != nil {
		return models.User{}, "", time.Time{}, err
	}
	return user, signed, exp, nil
}

func (s *AuthService) Current(ctx context.Context, userID string) (models.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.User{}, apperr.ErrUnauthorized
	}
	return user, err
}
