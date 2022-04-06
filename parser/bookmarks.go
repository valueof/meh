package parser

import (
	"io"

	"github.com/valueof/mold/util"
	"golang.org/x/net/html"
)

func ParseBookmarks(dat io.Reader) ([]Post, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	bookmarks := []Post{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			p := Post{}
			for t := n.FirstChild; t != nil; t = t.NextSibling {
				if t.Type != html.ElementNode {
					continue
				}

				switch {
				case t.Data == "a" && util.GetNodeAttr(t, "class") == "h-cite":
					p.Url = util.GetNodeAttr(t, "href")
					p.Id = util.ParseMediumId(p.Url)
					p.Title = util.GetNodeText(t)
				case t.Data == "time" && util.GetNodeAttr(t, "class") == "dt-published":
					p.PublishedAt = util.GetNodeText(t)
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
