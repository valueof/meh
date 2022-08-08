package server

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/valueof/meh/util"
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

	rand.Seed(time.Now().UnixNano())

	receipt := util.GenerateReceiptNumber()
	dest := filepath.Join(INBOUND_DIR, receipt)
	err = os.Mkdir(dest, 0700)

	for err != nil {
		if errors.Is(err, os.ErrExist) {
			logger.Printf("Receipt number collision, need to generate a new one")
			receipt = util.GenerateReceiptNumber()
			dest = filepath.Join(INBOUND_DIR, receipt)
			err = os.Mkdir(dest, 0700)
			continue
		}

		logger.Printf("Failed to create holding directory %s: %v\n", dest, err)
		internalServerError(w, r)
		return
	}

	dest = filepath.Join(dest, "upload.zip")
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
	go unzipAndParse(receipt, withImages, logger)

	url := fmt.Sprintf("/result/%s", receipt)
	http.Redirect(w, r, url, http.StatusFound)
}

func result(w http.ResponseWriter, r *http.Request) {
	logger := getLoggerFromContext(r.Context())

	receipt := strings.TrimPrefix(r.URL.Path, "/result/")
	if receipt == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if r.URL.Query().Has("dl") {
		file, err := os.Open(filepath.Join(INBOUND_DIR, receipt, "output.zip"))
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

		go cleanup(receipt, logger)
		return
	}

	st, exists := tasks.Status(receipt)
	if !exists {
		notFound(w, r)
		return
	}

	switch st {
	case TaskDone:
		render(w, r, "fetch.html", pageMeta{
			Title:      "[meh] Downloading...",
			SkipFooter: true,
			Refresh:    fmt.Sprintf("0;url=/result/%s/?dl", receipt),
		})
	case TaskErrUnknown:
		internalServerError(w, r)
	case TaskErrZipFormat:
		serverError(w, r, "The file we received wasn’t a valid zip file")
		go cleanup(receipt, logger)
	case TaskErrArchiveFormat:
		serverError(w, r, "The file we received wasn’t a valid Medium archive")
		go cleanup(receipt, logger)
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
