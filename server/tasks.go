package server

import (
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
	TASK_RUNNING taskStatus = 1
	TASK_DONE    taskStatus = 2
	TASK_ERROR   taskStatus = 3
)

type TaskPool struct {
	mu   sync.Mutex
	pool map[string]taskStatus
}

func (t *TaskPool) Create(h string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.pool[h] = TASK_RUNNING
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
		t.pool[h] = TASK_DONE
		return nil
	}

	return errors.New("can't complete task that doesn't exist")
}

func (t *TaskPool) Error(h string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.pool[h]; ok {
		t.pool[h] = TASK_ERROR
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
		tasks.Error(h)
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
		tasks.Error(h)
		return
	}

	output := filepath.Join(INBOUND_DIR, h, ".output")
	w := formatters.NewJSONFormatter(output, *logger)
	p := parser.NewParser(input, *logger, w)
	err = p.Parse()
	if err != nil {
		logger.Printf("parser.Parse(): %v", err)
		tasks.Error(h)
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
		tasks.Error(h)
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
