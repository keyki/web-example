package web

import (
    "log"
    "net/http"
)

type LoggingMiddleware struct {
    handler http.Handler
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s", r.Method, r.URL.Path)
    l.handler.ServeHTTP(w, r)
}

func newMiddleware(handler http.Handler) *LoggingMiddleware {
    return &LoggingMiddleware{handler: handler}
}
