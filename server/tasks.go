package server

import (
	"archive/zip"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/valueof/meh/formatters"
	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/util"
)

type taskStatus int

const (
	TaskRunning      taskStatus = 1
	TaskDone         taskStatus = 2
	TaskErrUnknown   taskStatus = 3
	TaskErrZipFormat taskStatus = 4
)

type TaskPool struct {
	mu   sync.Mutex
	pool map[string]taskStatus
}

func (t *TaskPool) Create(h string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.pool[h] = TaskRunning
}

func (t *TaskPool) Status(h string) (taskStatus, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if v, ok := t.pool[h]; ok {
		return v, true
	}

	return 0, false
}

func (t *TaskPool) Complete(h string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.pool[h]; ok {
		t.pool[h] = TaskDone
		return nil
	}

	return errors.New("can't complete task that doesn't exist")
}

func (t *TaskPool) Error(h string, e error) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.pool[h]; ok {
		switch e {
		case zip.ErrFormat:
			t.pool[h] = TaskErrZipFormat
		default:
			t.pool[h] = TaskErrUnknown
		}

		return nil
	}

	return errors.New("can't error task that doesn't exist")
}

func unzipAndParse(h string, withImages bool, logger *log.Logger) {
	tasks.Create(h)

	zip := filepath.Join(INBOUND_DIR, h, "upload.zip")
	tmp := filepath.Join(INBOUND_DIR, h, ".upload")
	err := util.UnzipArchive(zip, tmp)
	if err != nil {
		logger.Printf("UnzipArchive(%s, %s): %v", zip, tmp, err)
		tasks.Error(h, err)
		return
	}

	defer func() {
		logger.Printf("clean up: removing %s", zip)
		os.RemoveAll(zip)

		logger.Printf("clean up: removing %s", tmp)
		os.RemoveAll(tmp)
	}()

	input := util.FindArchiveRoot(tmp)
	input, err = filepath.Abs(input)
	fmt.Println(input)
	if err != nil {
		logger.Printf("filepath.Abs(): %v", err)
		tasks.Error(h, err)
		return
	}

	output := filepath.Join(INBOUND_DIR, h, ".output")
	w := formatters.NewJSONFormatter(output, *logger)
	p := parser.NewParser(input, *logger, w)
	err = p.Parse()
	if err != nil {
		logger.Printf("parser.Parse(): %v", err)
		tasks.Error(h, err)
		return
	}

	if withImages {
		p.FetchImages(output)
	}

	defer func() {
		logger.Printf("clean up: removing %s", output)
		os.RemoveAll(output)
	}()

	outzip := filepath.Join(INBOUND_DIR, h, "output.zip")
	err = util.ZipArchive(output, outzip)
	if err != nil {
		logger.Printf("util.ZipArchive(%s, %s): %v", output, outzip, err)
		tasks.Error(h, err)
		return
	}

	tasks.Complete(h)
}

func cleanup(h string, logger *log.Logger) {
	dir := filepath.Join(INBOUND_DIR, h)
	logger.Printf("Cleaning up %s", dir)

	err := os.RemoveAll(dir)
	if err != nil {
		logger.Printf("os.RemoveAll returned error: %v\n", err)
	}
}
