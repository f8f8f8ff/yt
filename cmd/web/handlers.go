package main

import (
	"html/template"
	"net/http"
	"net/url"
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
	requrl := "https://www.youtube.com/watch?v=sCNj0WMBkrs"
	path, err := app.dlmanager.GetVideo(requrl)
	if err != nil {
		app.serverError(w, err)
	}
	path = url.PathEscape("/dl/" + path)
	http.Redirect(w, r, path, http.StatusSeeOther)
}
