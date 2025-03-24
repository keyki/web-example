package web

import (
    "encoding/base64"
    "errors"
    "log"
    "net/http"
    "strings"
    "time"
    "web-example/types"
    "web-example/user"
    "web-example/util"
)

type Middleware func(user.UserRepository, http.Handler) http.Handler

func CreateStack(middlewares ...Middleware) Middleware {
    return func(userRepo user.UserRepository, next http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](userRepo, next)
        }
        return next
    }
}

func AuthenticationMiddleware(userRepo user.UserRepository, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Checking authentication")

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            log.Println("No Authorization header")
            util.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
            return
        }

        authParts := strings.SplitN(authHeader, " ", 2)
        if len(authParts) != 2 || authParts[0] != "Basic" {
            util.WriteError(w, http.StatusUnauthorized, errors.New("Only Basic auth is supported"))
            return
        }

        decoded, err := base64.StdEncoding.DecodeString(authParts[1])
        if err != nil {
            log.Printf("Error decoding base64 auth: %v", err)
            util.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
            return
        }

        creds := strings.SplitN(string(decoded), ":", 2)
        if len(creds) != 2 {
            util.WriteError(w, http.StatusUnauthorized, errors.New("Username format is incorrect"))
        }

        username := creds[0]
        userFromDb, err := userRepo.FindByUsername(username)
        if err != nil {
            log.Printf("Error finding user: %v", err)
            util.WriteError(w, http.StatusUnauthorized, errors.New("User not found"))
            return
        }
        if !util.CheckPassword(userFromDb.Password, creds[1]) {
            util.WriteError(w, http.StatusUnauthorized, errors.New("Incorrect password"))
            return
        }
        if types.ADMIN != userFromDb.Role && r.Method == "POST" {
            util.WriteError(w, http.StatusUnauthorized, errors.New("Admin role is required"))
            return
        }

        next.ServeHTTP(w, r)
    })
}

func MeasureMiddleware(_ user.UserRepository, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s took %v", r.Method, r.URL.Path, time.Since(startTime))
    })
}
