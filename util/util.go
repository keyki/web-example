package util

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"web-example/log"
	"web-example/types"
)

func DecodeJSON[T any](r *http.Request) (T, error) {
	var result T
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	code := 200
	if status != 0 {
		code = status
	}
	if err := writeJSON(w, code, v); err != nil {
		log.BaseLogger().Errorf("WriteJSON Error: %v", err)
		WriteError(w, http.StatusInternalServerError, err)
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func IsValidRole(role string) bool {
	valid := false
	for _, r := range types.Roles {
		if string(r) == role {
			valid = true
		}
	}
	return valid
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func NewInternalError() error {
	return errors.New("Internal error happened, please try again later.")
}

func GetUsername(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

func SetReqID(ctx context.Context) context.Context {
	return context.WithValue(ctx, types.ContextKeyReqID, uuid.New().String())
}

func GetReqID(ctx context.Context) string {
	return ctx.Value(types.ContextKeyReqID).(string)
}
