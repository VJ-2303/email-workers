package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheclHandler)
	mux.HandleFunc("/v1/send", app.sendEmailHandler)

	return app.limit(app.loggingMiddleware(mux))
}
