package web

import (
    "errors"
    "fmt"
    "github.com/jackc/pgx/v5/pgconn"
    "gorm.io/gorm"
    "log"
    "net/http"
    "web-example/user"
)

type Server struct {
    port int
    db   *gorm.DB
}

func NewApiServer(port int, db *gorm.DB) *Server {
    return &Server{port: port, db: db}
}

func (s *Server) Listen() {
    log.Printf("Listening on port %d", s.port)

    userStore := user.NewUserStore(s.db)
    userHandler := user.NewUserHandler(userStore)

    initAdminUser(userStore)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /users", userHandler.ListAll)
    //mux.HandleFunc("POST /user", createUser)

    err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), newMiddleware(mux))
    if err != nil {
        log.Fatalf("Failed to listen on port %d: %v", s.port, err)
    }
}

func initAdminUser(userStore *user.UserStore) {
    err := userStore.Create(user.User{
        UserName: "admin",
        Role:     user.ADMIN,
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
