package main

import (
	"html/template"
	"net/http"
	url_ "net/url"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.download(w, r, r.URL.RequestURI(), "audio")
		return
	}
	files := []string{
		"./ui/html/home.tmpl",
		"./ui/html/form.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) get(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/get" {
		app.getPost(w, r)
		return
	}
	app.clientError(w, http.StatusMethodNotAllowed)
}

func (app *application) getPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	url := r.PostForm.Get("url")
	medium := r.PostForm.Get("medium")
	app.download(w, r, url, medium)
}

func (app *application) download(w http.ResponseWriter, r *http.Request, url, medium string) {
	app.infoLog.Printf("url: %v, medium: %v", url, medium)
	var (
		path string
		err  error
	)
	path, err = app.dlmanager.Get(url, medium)
	if err != nil {
		app.serverError(w, err)
		app.clientErrorV(w, http.StatusBadRequest, err)
		return
	}
	path = url_.PathEscape("/dl/" + path)
	http.Redirect(w, r, path, http.StatusSeeOther)
}
