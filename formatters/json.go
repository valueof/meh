package formatters

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// JSONFormatter converts export data into raw JSON. It doesn't
// distinguish between different types of data. It creates a new
// .json file per each invokation of WriteFile.
type JSONFormatter struct {
	logger log.Logger
	root   string
}

func NewJSONFormatter(root string, logger log.Logger) *JSONFormatter {
	return &JSONFormatter{
		logger: logger,
		root:   root,
	}
}

func (w *JSONFormatter) WriteFile(fp string, v any) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		w.logger.Printf("can't marshal output for %s", fp)
		return err
	}

	// Make sure all directories exist to host this file
	dir := filepath.Dir(filepath.Join(w.root, fp+".json"))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		w.logger.Printf("can't create %s", dir)
		return err
	}

	dest := filepath.Join(w.root, fp+".json")
	err = os.WriteFile(dest, out, 0644)
	if err != nil {
		w.logger.Printf("can't write to %s", dest)
		return err
	}

	return nil
}
