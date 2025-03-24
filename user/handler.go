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

type Repository interface {
    ListAll() ([]*User, error)
    Create(user *User) error
    FindByUsername(username string) (*User, error)
}

type Handler struct {
    store Repository
}

func NewHandler(store Repository) *Handler {
    initAdminUser(store)
    return &Handler{store: store}
}

func convertToUserResponse(users []*User) (r []*Response) {
    for _, u := range users {
        r = append(r, u.ToReponse())
    }
    return r
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
    users, err := h.store.ListAll()
    if err != nil {
        util.WriteError(w, http.StatusInternalServerError, err)
    }
    util.WriteJSON(w, http.StatusOK, convertToUserResponse(users))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var userRequest Request
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    userName := r.PathValue("userName")
    log.Printf("Find user %s\n", userName)
    if userName == "" {
        util.WriteError(w, http.StatusBadRequest, errors.New("UserName is required"))
        return
    }

    user, err := h.store.FindByUsername(userName)
    if err != nil {
        log.Printf("Find error: %v", err)
        util.WriteJSON(w, http.StatusNotFound, []*Response{})
        return
    }

    util.WriteJSON(w, http.StatusOK, user.ToReponse())

}

func initAdminUser(store Repository) {
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
