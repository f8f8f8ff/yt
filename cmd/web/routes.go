package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(app.dlmanager.Downloader.Dir))
	mux.Handle("/dl/", app.logRequest(http.StripPrefix("/dl", fs)))

	mux.Handle("/", app.logRequest(http.HandlerFunc(app.home)))
	mux.Handle("/get", app.logRequest(http.HandlerFunc(app.get)))

	return mux
}
