package user

import (
    "encoding/json"
    "errors"
    "github.com/jackc/pgx/v5/pgconn"
    "log"
    "net/http"
    "web-example/types"
    "web-example/util"
)

type UserHandler struct {
    store UserRepository
}

func NewUserHandler(store UserRepository) *UserHandler {
    initAdminUser(store)
    return &UserHandler{store: store}
}

func initAdminUser(store UserRepository) {
    err := store.Create(&User{
        UserName: "admin",
        Role:     types.ADMIN,
    })
    if err == nil {
        return
    }
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        if pgErr.Code == "23505" {
            log.Printf("Admin user already exists")
        } else {
            log.Println(err)
        }
    } else {
        log.Printf("Cannot create admin user: %v", err)
    }
}

func (h *UserHandler) ListAll(w http.ResponseWriter, r *http.Request) {
    users, err := h.store.ListAll()
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, err)
    }
    util.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        util.WriteError(w, http.StatusBadRequest, err)
        return
    }
    if err := user.Validate(); err != nil {
        util.WriteError(w, http.StatusBadRequest, err)
        return
    }
    if err := h.store.Create(&user); err != nil {
        util.WriteError(w, http.StatusInternalServerError, err)
        log.Printf("Create Error: %v", err)
    }
}
