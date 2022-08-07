package util

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestParser(d string, t *testing.T, fn func(io.Reader, io.Reader) bool) {
	files, _ := ioutil.ReadDir(d)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ".html") {
			in, err := os.Open(path.Join(d, f.Name()))
			if err != nil {
				t.Errorf("can't open %s", f.Name())
			}
			defer in.Close()

			fout := strings.ReplaceAll(f.Name(), ".html", ".json")
			out, err := os.Open(path.Join(d, fout))
			if err != nil {
				t.Errorf("can't open %s", fout)
			}
			defer out.Close()

			if !fn(in, out) {
				t.Errorf("%s is unmatched with %s", f.Name(), fout)
			}
		}
	}
}
