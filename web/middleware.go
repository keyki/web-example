package web

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
	"web-example/log"
	"web-example/types"
	"web-example/user"
	"web-example/util"
)

type Middleware func(user.Repository, http.Handler) http.Handler

func CreateMiddleware(middlewares ...Middleware) Middleware {
	return func(userRepo user.Repository, next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](userRepo, next)
		}
		return next
	}
}

func AuthenticationMiddleware(userRepo user.Repository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Logger(r.Context()).Info("Checking authentication")

		username, password, ok := r.BasicAuth()
		if !ok {
			log.Logger(r.Context()).Info("No Authorization header")
			util.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		userFromDb, err := userRepo.FindByUsername(r.Context(), username)
		if err != nil {
			log.Logger(r.Context()).Infof("Error finding user: %v", err)
			util.WriteError(w, http.StatusUnauthorized, errors.New("User not found"))
			return
		}
		if !util.CheckPassword(userFromDb.Password, password) {
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

func MeasureMiddleware(_ user.Repository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		log.Logger(r.Context()).Infof("%s %s took %v", r.Method, r.URL.Path, time.Since(startTime))
	})
}

func RequestIdMiddleware(_ user.Repository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		ctx := context.WithValue(r.Context(), types.ContextKeyReqID, requestId)
		r = r.WithContext(ctx)
		w.Header().Add(types.HTTPHeaderRequestID, requestId)
		next.ServeHTTP(w, r)
	})
}
