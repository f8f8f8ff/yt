package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/home.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) video(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		w.Header().Set("Allow", http.MethodPost)
	// 		app.clientError(w, http.StatusMethodNotAllowed)
	// 		return
	// 	}
	w.Write([]byte("video"))
}
