package util

import (
    "encoding/json"
    "log"
    "net/http"
    "web-example/types"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
    if err := writeJSON(w, 200, v); err != nil {
        log.Printf("WriteJSON Error: %v", err)
        WriteError(w, http.StatusInternalServerError, err)
    }
}

func WriteError(w http.ResponseWriter, status int, err error) {
    if err := writeJSON(w, status, map[string]string{"error": err.Error()}); err != nil {
        log.Printf("WriteError Error: %v", err)
    }
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
