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

    userMux := http.NewServeMux()
    userMux.HandleFunc("GET /users", userHandler.ListAll)
    userMux.HandleFunc("POST /user", userHandler.Create)

    userMiddlewareStack := CreateStack(AuthenticationMiddleware)
    wrappedUserMux := userMiddlewareStack(userStore, http.StripPrefix("/v1", userMux))

    mainMux := http.NewServeMux()
    mainMux.HandleFunc("GET /info", Info)
    mainMux.Handle("/v1/", wrappedUserMux)

    mainMiddlewareStack := CreateStack(MeasureMiddleware)
    wrappedMainMux := mainMiddlewareStack(nil, mainMux)

    err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), wrappedMainMux)
    if err != nil {
        log.Fatalf("Failed to listen on port %d: %v", s.port, err)
    }
}

func Info(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK")
}
