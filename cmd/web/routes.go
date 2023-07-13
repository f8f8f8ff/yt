package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", app.logRequest(http.HandlerFunc(app.home)))
	mux.Handle("/get", app.logRequest(http.HandlerFunc(app.get)))
	mux.Handle("/dl/", app.logRequest(http.HandlerFunc(app.dl)))
	mux.Handle("/zip", app.logRequest(http.HandlerFunc(app.zip)))

	fs := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
