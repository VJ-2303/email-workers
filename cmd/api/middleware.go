package main

import (
	"net/http"
	"time"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		app.logger.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
