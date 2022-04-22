package parser

import (
	"io"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func ParseBookmarks(dat io.Reader) ([]schema.Post, error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	bookmarks := []schema.Post{}

	var f func(*util.Node)
	f = func(n *util.Node) {
		if n.IsElement("li") {
			p := schema.Post{}
			for t := n.FirstChild; t != nil; t = t.NextSibling {
				if t.Type != html.ElementNode {
					continue
				}

				switch {
				case t.Data == "a" && t.Attrs["class"] == "h-cite":
					p.Url = t.Attrs["href"]
					p.Id = util.ParseMediumId(p.Url)
					p.Title = t.Text()
				case t.Data == "time" && t.Attrs["class"] == "dt-published":
					p.PublishedAt = t.Text()
				}
			}
			bookmarks = append(bookmarks, p)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return bookmarks, nil
}
