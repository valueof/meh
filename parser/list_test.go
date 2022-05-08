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

func TestParseLists(t *testing.T) {
	util.TestParser("../testdata/lists", t, func(in, out io.Reader) bool {
		have, err := parser.ParseList(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want schema.List
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		fmt.Println(have)

		return reflect.DeepEqual(have, &want)
	})
}
