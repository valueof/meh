package parser_test

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/schema"
)

func TestParseIPs(t *testing.T) {
	tests := map[string][]schema.IP{
		"ips/ips.html": {
			{
				Address:   "127.0.0.1",
				CreatedAt: "2022-03-16 4:39 am",
			},
			{
				Address:   "2001:db8::68",
				CreatedAt: "2022-04-16 4:39 am",
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata/", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		ips, err := parser.ParseIps(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := ips[i]

			if reflect.DeepEqual(have, want) == false {
				t.Errorf("want: %v; have: %v", want, have)
			}
		}
	}
}
