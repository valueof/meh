package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var INBOUND_DIR string = "/Users/anton/Code/meh/tmp"

//go:embed html
var templates embed.FS

type pageMeta struct {
	Title string
}

func homepage(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	tmpl := template.Must(template.New("base.html").ParseFS(templates, "html/base.html", "html/home.html"))
	err := tmpl.Execute(w, pageMeta{Title: "Medium Export Helper"})
	if err != nil {
		panic("Can't execute template")
	}
}

func convert(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	uploads := req.MultipartForm.File["archive"]
	if len(uploads) == 0 {
		fmt.Fprintf(w, "no file was sent from the client")
		return
	}

	if len(uploads) > 1 {
		fmt.Fprintf(w, "too many files were sent from the client")
		return
	}

	header := uploads[0]
	file, err := header.Open()
	if err != nil {
		fmt.Fprintf(w, "FileHeader.Open() err: %v", err)
		return
	}

	defer file.Close()

	dest, err := os.Create(path.Join(INBOUND_DIR, header.Filename))
	if err != nil {
		fmt.Fprintf(w, "Couldn't create dest file: %v", err)
		return
	}

	defer dest.Close()
	_, err = io.Copy(dest, file)
	if err != nil {
		fmt.Fprintf(w, "io.Copy err: %v", err)
		return
	}

	fmt.Fprintf(w, "success! %s was uploaded", header.Filename)
}

func RunHTTPServer(addr string) (s *http.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/convert", convert)

	s = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         addr,
		Handler:      mux,
	}

	log.Print(s.ListenAndServe())
	return
}
