package parser_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func TestParsePost(t *testing.T) {
	util.TestParser("../testdata/posts", t, func(in, out io.Reader) bool {
		have, err := parser.ParsePost(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want schema.Post
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		return reflect.DeepEqual(have, &want)
	})
}
