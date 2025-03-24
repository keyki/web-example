package web

import (
    "fmt"
    "gorm.io/gorm"
    "log"
    "net/http"
    "web-example/types"
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

    mux := http.NewServeMux()
    mux.HandleFunc("GET /info", Info)
    mux.HandleFunc("GET /users", Authenticator(userHandler.ListAll, userStore, types.USER))
    mux.HandleFunc("POST /user", Authenticator(userHandler.Create, userStore, types.ADMIN))

    middleware := CreateStack(
        MeasureMiddleware,
    )

    err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), middleware(mux))
    if err != nil {
        log.Fatalf("Failed to listen on port %d: %v", s.port, err)
    }
}

func Info(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK")
}
