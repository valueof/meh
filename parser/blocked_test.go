package parser_test

import (
	"os"
	"path"
	"testing"

	"github.com/valueof/meh/parser"
)

func TestParseBlocked(t *testing.T) {
	tests := map[string][]parser.User{
		"blocks/blocked.html": {
			{Username: "bob"},
			{Username: "lelandpalmer"},
			{Username: "leojohnson"},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata/", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		blocked, err := parser.ParseBlocked(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i := 0; i < len(tt); i++ {
			want := tt[i].Username
			have := blocked[i].Username
			if want != have {
				t.Errorf("want: %s; have: %s", want, have)
			}
		}
	}
}
