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

func convertToUserResponse(users []*User) (r []*UserResponse) {
    for _, u := range users {
        r = append(r, u.ToReponse())
    }
    return r
}

func (h *UserHandler) ListAll(w http.ResponseWriter, r *http.Request) {
    users, err := h.store.ListAll()
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, err)
    }
    util.WriteJSON(w, http.StatusOK, convertToUserResponse(users))
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var userRequest UserRequest
    if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
        util.WriteError(w, http.StatusBadRequest, err)
        return
    }
    if err := userRequest.Validate(); err != nil {
        util.WriteError(w, http.StatusBadRequest, err)
        return
    }
    userRequest.Password = util.HashPassword(userRequest.Password)
    if err := h.store.Create(userRequest.ToUser()); err != nil {
        util.WriteError(w, http.StatusInternalServerError, err)
        log.Printf("Create Error: %v", err)
    }
}

func initAdminUser(store UserRepository) {
    err := store.Create(&User{
        UserName: "admin",
        Password: util.HashPassword("admin"),
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
