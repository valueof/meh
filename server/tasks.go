package server

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type taskStatus int

const (
	TASK_RUNNING taskStatus = 1
	TASK_DONE    taskStatus = 2
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

func unzipAndParse(h string, logger *log.Logger) {
	tasks.Create(h)

	// TODO:
	// - unzip into ./input
	// - parse util.FindArchiveRoot(./input)
	// - zip into ./output.zip
	time.Sleep(15 * time.Second)

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
