package user

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type UserHandler struct {
    store UserRepository
}

func NewUserHandler(store UserRepository) *UserHandler {
    return &UserHandler{store: store}
}

func (h *UserHandler) ListAll(w http.ResponseWriter, r *http.Request) {
    users, err := h.store.ListAll()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        _, err := fmt.Fprintf(w, "Failed to list users: %v", err)
        if err != nil {
            log.Printf("Failed to send error response: %v", err)
        }
    }
    marshal, err := json.Marshal(users)
    if err != nil {
        log.Printf("Failed to marshal users: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        _, err = fmt.Fprintf(w, "Failed to send error response: %v", err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(marshal)
}
