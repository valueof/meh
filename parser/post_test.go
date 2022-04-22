package parser_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/schema"
)

func TestParsePost(t *testing.T) {
	dir := "../testdata/posts"
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ".html") {
			input, err := os.Open(path.Join(dir, f.Name()))
			defer input.Close()
			if err != nil {
				t.Errorf("can't open %s", f.Name())
			}

			fout := strings.ReplaceAll(f.Name(), ".html", ".json")
			output, err := os.Open(path.Join(dir, fout))
			defer output.Close()
			if err != nil {
				t.Errorf("can't open %s", fout)
			}

			have, err := parser.ParsePost(input)
			if err != nil {
				t.Errorf("error parsing input: %v", err)
			}

			var want schema.Post
			outb, _ := ioutil.ReadAll(output)
			json.Unmarshal(outb, &want)

			if reflect.DeepEqual(have, &want) == false {
				t.Errorf("%s is unmatched with %s", f.Name(), fout)
			}
		}
	}
}
