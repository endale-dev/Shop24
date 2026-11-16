package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"shop24/internal/apperr"
)

func Data(w http.ResponseWriter, status int, data any) {
	write(w, status, map[string]any{"data": data})
}

func Error(w http.ResponseWriter, err error) {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		write(w, appErr.Status, map[string]any{
			"error": map[string]string{"code": appErr.Code, "message": appErr.Message},
		})
		return
	}
	write(w, http.StatusInternalServerError, map[string]any{
		"error": map[string]string{"code": "internal", "message": "Something went wrong"},
	})
}

func write(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
