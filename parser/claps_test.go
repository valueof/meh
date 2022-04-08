package parser_test

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/valueof/meh/parser"
)

func TestParseClaps(t *testing.T) {
	tests := map[string][]parser.Clap{
		"claps/claps.html": {
			{
				Amount: 1,
				Post: parser.Post{
					Id:          "9e53ca408c48",
					Url:         "https://medium.com/p/welcome-to-medium-9e53ca408c48",
					Title:       "Welcome to Medium",
					PublishedAt: "2012-08-14 11:54 pm",
				},
			},
			{
				Amount: 25,
				Post: parser.Post{
					Id:          "b8d43e4c204d",
					Url:         "https://medium.com/p/re-thinking-j-school-b8d43e4c204d",
					Title:       "Re-thinking J-school",
					PublishedAt: "2013-04-28 6:55 am",
				},
			},
			{
				Amount: 50,
				Post: parser.Post{
					Id:          "3d26424537aa",
					Url:         "https://medium.com/p/i-accidentally-bought-a-banksy-in-2003-3d26424537aa",
					Title:       "I Accidentally Bought a Banksy in 2003",
					PublishedAt: "2013-05-5 10:08 pm",
				},
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata/", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		claps, err := parser.ParseClaps(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := claps[i]

			if reflect.DeepEqual(have, want) == false {
				t.Errorf("want: %v;\nhave: %v", want, have)
			}
		}
	}
}
