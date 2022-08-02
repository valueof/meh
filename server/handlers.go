package server

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type pageMeta struct {
	Title      string
	SkipFooter bool
	Refresh    string
}

type errorPageData struct {
	RequestID    string
	ErrorMessage string
	pageMeta
}

func notFound(w http.ResponseWriter, r *http.Request) {
	data := pageMeta{}
	data.Title = "[meh] Page Not Found"
	data.SkipFooter = true

	render(w, r, "404.html", data)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	data := errorPageData{}
	data.Title = "[meh] Internal Server Error"
	data.SkipFooter = true
	data.RequestID = getRequestIDFromContext(r.Context())

	render(w, r, "500.html", data)
}

func serverError(w http.ResponseWriter, r *http.Request, m string) {
	data := errorPageData{}
	data.Title = "[meh] Something went wrong"
	data.SkipFooter = true
	data.RequestID = getRequestIDFromContext(r.Context())
	data.ErrorMessage = m

	render(w, r, "500.html", data)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w, r)
		return
	}

	render(w, r, "home.html", pageMeta{Title: "Medium Export Helper"})
}

func upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := getLoggerFromContext(ctx)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		logger.Printf("ParseForm() err: %v", err)
		internalServerError(w, r)
		return
	}

	withImages := len(r.MultipartForm.Value["withImages"]) > 0
	uploads := r.MultipartForm.File["archive"]
	if len(uploads) == 0 {
		logger.Printf("no file was sent from the client")
		internalServerError(w, r)
		return
	}

	if len(uploads) > 1 {
		logger.Printf("too many files were sent from the client")
		internalServerError(w, r)
		return
	}

	header := uploads[0]
	file, err := header.Open()
	if err != nil {
		logger.Printf("FileHeader.Open() err: %v", err)
		internalServerError(w, r)
		return
	}

	defer file.Close()

	h := sha256.New()
	io.TeeReader(file, h)
	hashsum := fmt.Sprintf("%x", h.Sum(nil))

	dest := filepath.Join(INBOUND_DIR, hashsum)
	err = os.Mkdir(dest, 0700)
	if err != nil && !errors.Is(err, os.ErrExist) {
		logger.Printf("Failed to create holding directory %s: %v\n", dest, err)
		internalServerError(w, r)
		return
	}

	dest = filepath.Join(dest, "upload.zip")
	_, err = os.Stat(dest)
	if err == nil {
		// File already exists, check whether we need to reprocess it and redirect
		if t, ok := tasks.Status(hashsum); !ok && t != TaskRunning {
			go unzipAndParse(hashsum, withImages, logger)
		}

		url := fmt.Sprintf("/result/%s", hashsum)
		http.Redirect(w, r, url, http.StatusFound)
	} else if errors.Is(err, os.ErrNotExist) {
		// File doesn't exist, upload and send for processing
		upload, err := os.Create(dest)
		if err != nil {
			logger.Printf("Couldn't create dest file: %v", err)
			internalServerError(w, r)
			return
		}

		defer upload.Close()
		_, err = io.Copy(upload, file)
		if err != nil {
			logger.Printf("io.Copy err: %v", err)
			internalServerError(w, r)
			return
		}

		logger.Printf("Uploaded %s", dest)

		go unzipAndParse(hashsum, withImages, logger)

		url := fmt.Sprintf("/result/%s", hashsum)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// Some other error, run around in panic
		logger.Printf("os.Stat returned an unexpected error: %v\n", err)
		internalServerError(w, r)
	}
}

func result(w http.ResponseWriter, r *http.Request) {
	logger := getLoggerFromContext(r.Context())

	hashsum := strings.TrimPrefix(r.URL.Path, "/result/")
	if hashsum == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if r.URL.Query().Has("dl") {
		file, err := os.Open(filepath.Join(INBOUND_DIR, hashsum, "output.zip"))
		if err != nil {
			logger.Printf("Couldn't read file for download: %v\n", err)
			notFound(w, r)
			return
		}

		info, err := file.Stat()
		if err != nil {
			logger.Printf("file.Stat() returned an error: %v\n", err)
			internalServerError(w, r)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
		_, err = io.Copy(w, file)
		if err != nil {
			logger.Printf("io.Copy err: %v", err)
			internalServerError(w, r)
			return
		}

		go cleanup(hashsum, logger)
		return
	}

	st, exists := tasks.Status(hashsum)
	if !exists {
		notFound(w, r)
		return
	}

	switch st {
	case TaskDone:
		render(w, r, "fetch.html", pageMeta{
			Title:      "[meh] Downloading...",
			SkipFooter: true,
			Refresh:    fmt.Sprintf("0;url=/result/%s/?dl", hashsum),
		})
	case TaskErrUnknown:
		internalServerError(w, r)
	case TaskErrZipFormat:
		serverError(w, r, "The file we received wasn’t a valid zip file")
	default:
		render(w, r, "wait.html", pageMeta{
			Title:      "[meh] Converting...",
			SkipFooter: true,
			Refresh:    "10",
		})
	}
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}
