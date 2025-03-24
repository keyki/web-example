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

    userStore := user.NewStore(s.db)
    userHandler := user.NewHandler(userStore)

    userMux := http.NewServeMux()
    userMux.HandleFunc("GET /users", userHandler.ListAll)
    userMux.HandleFunc("POST /user", userHandler.Create)
    userMux.HandleFunc("GET /user/{userName}", userHandler.Get)

    userMiddleware := CreateMiddleware(AuthenticationMiddleware)
    wrappedUserMux := userMiddleware(userStore, http.StripPrefix("/v1", userMux))

    mainMux := http.NewServeMux()
    mainMux.HandleFunc("GET /info", Info)
    mainMux.Handle("/v1/", wrappedUserMux)

    mainMiddleware := CreateMiddleware(MeasureMiddleware)
    wrappedMainMux := mainMiddleware(nil, mainMux)

    err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), wrappedMainMux)
    if err != nil {
        log.Fatalf("Failed to listen on port %d: %v", s.port, err)
    }
}

func Info(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"version": "v1", "status": "ok"}`))
}
