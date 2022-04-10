package parser_test

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/valueof/meh/parser"
)

func TestParseInterestsPublications(t *testing.T) {
	tests := map[string][]parser.Publication{
		"interests/publications.html": {
			{
				Name: "Wildlife Trekker",
				Url:  "https://medium.com/wildlife-trekker",
			},
			{
				Name: "3 min read",
				Url:  "https://blog.medium.com",
			},
			{
				Name: "Writing by Dan Pupius",
				Url:  "https://writing.pupius.co.uk",
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		pubs, err := parser.ParseInterestsPublications(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := pubs[i]

			if reflect.DeepEqual(have, want) == false {
				t.Errorf("want: %v;\nhave: %v", want, have)
			}
		}
	}
}

func TestParseInterestsTags(t *testing.T) {
	tests := map[string][]parser.Tag{
		"interests/tags.html": {
			{
				Name: "Birds",
				Url:  "https://medium.com/tag/birds",
			},
			{
				Name: "Libraries",
				Url:  "https://medium.com/tag/libraries",
			},
			{
				Name: "Books",
				Url:  "https://medium.com/tag/books",
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		tags, err := parser.ParseInterestsTags(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := tags[i]

			if reflect.DeepEqual(have, want) == false {
				t.Errorf("want: %v;\nhave: %v", want, have)
			}
		}
	}
}

func TestParseInterestsTopics(t *testing.T) {
	tests := map[string][]parser.Topic{
		"interests/topics.html": {
			{
				Name: "Photography",
				Url:  "https://medium.com/topic/photography",
			},
			{
				Name: "Education",
				Url:  "https://medium.com/topic/education",
			},
			{
				Name: "Programming",
				Url:  "https://medium.com/topic/programming",
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		topics, err := parser.ParseInterestsTopics(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := topics[i]

			if reflect.DeepEqual(have, want) == false {
				t.Errorf("want: %v;\nhave: %v", want, have)
			}
		}
	}
}

func TestParseInterestsWriters(t *testing.T) {
	tests := map[string][]parser.User{
		"interests/writers.html": {
			{
				Name:     "Randy Runtsch",
				Username: "rruntsch",
				Url:      "https://medium.com/@rruntsch",
			},
			{
				Name:     "Aryannayakk",
				Username: "aryannayakk",
				Url:      "https://medium.com/@aryannayakk",
			},
			{
				Name:     "Dan Pupius",
				Username: "dpup",
				Url:      "https://medium.com/@dpup",
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		users, err := parser.ParseInterestsWriters(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := users[i]

			if reflect.DeepEqual(have, want) == false {
				t.Errorf("want: %v;\nhave: %v", want, have)
			}
		}
	}
}
