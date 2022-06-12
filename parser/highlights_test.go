package parser_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func TestParseHighlights(t *testing.T) {
	util.TestParser("../testdata/highlights", t, func(in, out io.Reader) bool {
		have, err := parser.ParseHighlights(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want []schema.Highlight
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		for i, v := range have {
			fmt.Println(v)
			fmt.Println(want[i])
		}

		return reflect.DeepEqual(have, want)
	})
}
