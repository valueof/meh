package server

import (
	"archive/zip"
	"errors"
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
	TaskRunning          taskStatus = 1
	TaskDone             taskStatus = 2
	TaskErrUnknown       taskStatus = 3
	TaskErrZipFormat     taskStatus = 4
	TaskErrArchiveFormat taskStatus = 5
)

type TaskPool struct {
	mu   sync.Mutex
	pool map[string]taskStatus
}

func (t *TaskPool) Create(receipt string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.pool[receipt] = TaskRunning
}

func (t *TaskPool) Status(receipt string) (taskStatus, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if v, ok := t.pool[receipt]; ok {
		return v, true
	}

	return 0, false
}

func (t *TaskPool) Complete(receipt string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.pool[receipt]; ok {
		t.pool[receipt] = TaskDone
		return nil
	}

	return errors.New("can't complete task that doesn't exist")
}

func (t *TaskPool) Error(receipt string, e error) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.pool[receipt]; ok {
		switch e {
		case zip.ErrFormat:
			t.pool[receipt] = TaskErrZipFormat
		case util.ErrArchiveRootNotFound:
			t.pool[receipt] = TaskErrArchiveFormat
		default:
			t.pool[receipt] = TaskErrUnknown
		}

		return nil
	}

	return errors.New("can't error task that doesn't exist")
}

func unzipAndParse(receipt string, withImages bool, logger *log.Logger) {
	tasks.Create(receipt)

	zip := filepath.Join(INBOUND_DIR, receipt, "upload.zip")
	tmp := filepath.Join(INBOUND_DIR, receipt, ".upload")
	err := util.UnzipArchive(zip, tmp)
	if err != nil {
		logger.Printf("UnzipArchive(%s, %s): %v", zip, tmp, err)
		tasks.Error(receipt, err)
		return
	}

	defer func() {
		logger.Printf("clean up: removing %s", zip)
		os.RemoveAll(zip)

		logger.Printf("clean up: removing %s", tmp)
		os.RemoveAll(tmp)
	}()

	root, err := util.FindArchiveRoot(tmp)
	if err != nil {
		logger.Printf("util.FindArchiveRoot(%s): %v", tmp, err)
		tasks.Error(receipt, err)
		return
	}

	input, err := filepath.Abs(root)
	if err != nil {
		logger.Printf("filepath.Abs(): %v", err)
		tasks.Error(receipt, err)
		return
	}

	output := filepath.Join(INBOUND_DIR, receipt, ".output")
	w := formatters.NewJSONFormatter(output, *logger)
	p := parser.NewParser(input, *logger, w)
	err = p.Parse()
	if err != nil {
		logger.Printf("parser.Parse(): %v", err)
		tasks.Error(receipt, err)
		return
	}

	if withImages {
		p.FetchImages(output)
	}

	defer func() {
		logger.Printf("clean up: removing %s", output)
		os.RemoveAll(output)
	}()

	outzip := filepath.Join(INBOUND_DIR, receipt, "output.zip")
	err = util.ZipArchive(output, outzip)
	if err != nil {
		logger.Printf("util.ZipArchive(%s, %s): %v", output, outzip, err)
		tasks.Error(receipt, err)
		return
	}

	tasks.Complete(receipt)
}

func cleanup(receipt string, logger *log.Logger) {
	dir := filepath.Join(INBOUND_DIR, receipt)
	logger.Printf("Cleaning up %s", dir)

	err := os.RemoveAll(dir)
	if err != nil {
		logger.Printf("os.RemoveAll returned error: %v\n", err)
	}
}
