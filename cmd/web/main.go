package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"yt/internal/ytdl"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	dlmanager *ytdl.Manager
}

func main() {
	addr := flag.String("addr", ":4000", "http network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		dlmanager: ytdl.NewManager(),
	}
	// app.dlmanager.Downloader.Flags = []string{"--no-simulate"}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
