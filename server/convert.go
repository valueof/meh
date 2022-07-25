package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func convert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := getLoggerFromContext(ctx)
	rid := getRequestIDFromContext(ctx)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		logger.Printf("ParseForm() err: %v", err)
		internalServerError(w, r)
		return
	}

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

	dest := filepath.Join(INBOUND_DIR, rid)
	err = os.Mkdir(dest, 0700)
	if err != nil {
		logger.Printf("Failed to create holding directory %s: %v\n", dest, err)
		internalServerError(w, r)
		return
	}

	upload, err := os.Create(filepath.Join(dest, header.Filename))
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

	// TODO:
	//  - Generate and show a receipt number
	fmt.Fprintf(w, "success! %s was uploaded", header.Filename)
}
