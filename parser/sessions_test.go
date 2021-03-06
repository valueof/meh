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

func TestParseSessions(t *testing.T) {
	util.TestParser("../testdata/sessions", t, func(in, out io.Reader) bool {
		have, err := parser.ParseSessions(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want []schema.Session
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		return reflect.DeepEqual(have, want)
	})
}
