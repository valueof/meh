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

func TestParsePublicationsFollowing(t *testing.T) {
	util.TestParser("../testdata/following/publications", t, func(in, out io.Reader) bool {
		have, err := parser.ParsePublicationFollowing(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want []schema.Publication
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		return reflect.DeepEqual(have, want)
	})
}

func TestParseTopicsFollowing(t *testing.T) {
	util.TestParser("../testdata/following/topics", t, func(in, out io.Reader) bool {
		have, err := parser.ParseTopicsFollowing(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want []schema.Topic
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		return reflect.DeepEqual(have, want)
	})
}

func TestParseUsersFollowing(t *testing.T) {
	util.TestParser("../testdata/following/users", t, func(in, out io.Reader) bool {
		have, err := parser.ParseUsersFollowing(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want []schema.User
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		return reflect.DeepEqual(have, want)
	})
}

func TestParseUsersSuggested(t *testing.T) {
	util.TestParser("../testdata/following/suggested", t, func(in, out io.Reader) bool {
		have, err := parser.ParseUsersSuggested(in)
		if err != nil {
			t.Errorf("error parsing input: %v", err)
		}

		var want []schema.User
		outb, _ := ioutil.ReadAll(out)
		json.Unmarshal(outb, &want)

		return reflect.DeepEqual(have, want)
	})
}
