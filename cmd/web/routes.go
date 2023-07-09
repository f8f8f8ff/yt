package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/video", app.video)
	fs := http.FileServer(http.Dir(app.dlmanager.Downloader.Dir))
	mux.Handle("/dl", http.StripPrefix("/dl", fs))

	return mux
}
