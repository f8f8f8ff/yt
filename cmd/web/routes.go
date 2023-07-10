package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", app.logRequest(http.HandlerFunc(app.home)))
	mux.Handle("/get", app.logRequest(http.HandlerFunc(app.get)))
	mux.Handle("/dl/", app.logRequest(http.HandlerFunc(app.dl)))

	return mux
}
