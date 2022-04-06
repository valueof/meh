package parser_test

import (
	"os"
	"path"
	"testing"

	"github.com/valueof/meh/parser"
)

func TestParseBookmarks(t *testing.T) {
	tests := map[string][]parser.Post{
		"bookmarks/bookmarks.html": {
			{
				Id:          "e06382acd276",
				Url:         "https://medium.com/p/sprint-burndown-charts-gone-wrong-e06382acd276",
				Title:       "Sprint Burndown Charts Gone Wrong",
				PublishedAt: "2020-08-21 4:50 pm",
			},
			{
				Id:          "75aa4ed61c",
				Url:         "https://medium.com/p/dont-call-it-a-trend-a-brief-history-of-organizing-in-tech-75aa4ed61c",
				Title:       "Don’t Call It a Trend: A Brief History of Organizing in Tech",
				PublishedAt: "2020-02-19 2:54 am",
			},
			{
				Id:          "d8f306e18fbe",
				Url:         "https://medium.com/p/my-productivity-stopped-me-from-growing-d8f306e18fbe",
				Title:       "My Productivity Stopped Me From Growing",
				PublishedAt: "2019-12-14 8:41 pm",
			},
			{
				Id:          "e23a49a52696",
				Url:         "https://medium.com/p/life-of-a-vegetarian-in-japan-e23a49a52696",
				Title:       "Life of a Vegetarian in Japan",
				PublishedAt: "2019-06-20 5:16 am",
			},
		},
	}

	for fp, tt := range tests {
		dat, err := os.Open(path.Join("../testdata/", fp))
		if err != nil {
			t.Errorf("no testdata file: %s", fp)
			return
		}

		posts, err := parser.ParseBookmarks(dat)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		for i, want := range tt {
			have := posts[i]

			if want.Id != have.Id {
				t.Errorf("want: %s; have: %s", want.Id, have.Id)
			}

			if want.Url != have.Url {
				t.Errorf("want: %s; have: %s", want.Url, have.Url)
			}

			if want.Title != have.Title {
				t.Errorf("want: %s; have: %s", want.Title, have.Title)
			}

			if want.PublishedAt != have.PublishedAt {
				t.Errorf("want: %s; have: %s", want.PublishedAt, have.PublishedAt)
			}
		}
	}
}
