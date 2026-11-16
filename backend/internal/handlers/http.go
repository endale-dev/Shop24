package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"shop24/internal/apperr"
	"shop24/internal/httpx"
	"shop24/internal/token"
)

func decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return apperr.New(http.StatusRequestEntityTooLarge, "body_too_large", "Request body is too large")
		}
		return apperr.New(http.StatusBadRequest, "invalid_json", "Request body is not valid JSON")
	}
	return nil
}

func currentUser(r *http.Request) (string, error) {
	id, ok := token.UserID(r.Context())
	if !ok {
		return "", apperr.ErrUnauthorized
	}
	return id, nil
}

func fail(log *slog.Logger, w http.ResponseWriter, err error) {
	if _, ok := apperr.As(err); !ok {
		log.Error("request failed", "error", err)
	}
	httpx.Error(w, err)
}
