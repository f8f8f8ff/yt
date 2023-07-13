package main

import (
	"html/template"
	"net/http"
	_url "net/url"
	"path"
	_path "path"

	"yt/internal/dirlist"
	"yt/internal/dl"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		url, err := dl.PartialUrlToYoutube(r.URL.RequestURI())
		if err != nil {
			app.notFound(w)
			return
		}
		app.download(w, r, url, dl.MediumAudio, false)
		return
	}
	files := []string{
		"./ui/html/home.tmpl",
		"./ui/html/form.tmpl",
		"./ui/html/info.tmpl",
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
	medium := mediumFromString(r.PostForm.Get("medium"))
	app.download(w, r, url, medium, true)
}

func (app *application) download(w http.ResponseWriter, r *http.Request, url string, medium dl.Medium, redirect bool) {
	app.infoLog.Printf("url: %v, medium: %v", url, medium)
	var (
		paths []string
		err   error
	)
	paths, err = app.dlmanager.Get(url, medium)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if len(paths) > 1 {
		http.Redirect(w, r, "/dl", http.StatusSeeOther)
		return
	}
	dlPath := _url.PathEscape(paths[0])
	w.Header().Set("Content-Disposition", "attachment; filename="+dlPath)
	if redirect {
		w.Header().Set("Location", "/")
	}
	fp := _path.Join(app.dlmanager.Dir(), dlPath)
	http.ServeFile(w, r, fp)
}

func (app *application) dl(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dl/" {
		fs := http.StripPrefix("/dl", http.FileServer(http.Dir(app.dlmanager.Downloader.Dir)))
		fs.ServeHTTP(w, r)
		return
	}
	templatefiles := []string{
		"./ui/html/home.tmpl",
		"./ui/html/archive.tmpl",
	}
	ts, err := template.ParseFiles(templatefiles...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	files, err := dirlist.Dir(app.dlmanager.Dir())
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := struct{ Files []dirlist.File }{
		Files: files,
	}
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) zip(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/zip" {
		app.zipPost(w, r)
		return
	}
	app.clientError(w, http.StatusMethodNotAllowed)
}

func (app *application) zipPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	files := r.PostForm["files"]
	// TODO if single file is selected, don't zip
	if len(files) < 2 {
		http.Redirect(w, r, "/dl", http.StatusSeeOther)
		return
	}
	for i, f := range files {
		files[i], err = _url.PathUnescape(f)
		if err != nil {
			app.serverError(w, err)
		}
	}
	z, err := app.dlmanager.Zip(files...)
	if err != nil {
		app.serverError(w, err)
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(z))
	http.ServeFile(w, r, z)
}
