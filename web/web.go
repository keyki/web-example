package web

import (
    "fmt"
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

    mux := http.NewServeMux()
    mux.HandleFunc("GET /users", userHandler.ListAll)
    mux.HandleFunc("POST /user", userHandler.Create)

    middleware := CreateStack(
        Authentication,
        Measure,
    )

    err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), middleware(mux))
    if err != nil {
        log.Fatalf("Failed to listen on port %d: %v", s.port, err)
    }
}
