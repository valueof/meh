package parser

import (
	"io"

	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func parseFooter(n *html.Node, post *Post) {
	switch {
	case util.IsElement(n, "time") && util.HasClass(n, "dt-published"):
		post.PublishedAt = util.GetNodeAttr(n, "datetime")
	case util.IsElement(n, "a") && util.HasClass(n, "p-canonical"):
		post.Url = util.GetNodeAttr(n, "href")
		post.Id = util.ParseMediumId(post.Url)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseFooter(c, post)
	}
}

func ParsePost(dat io.Reader) (*Post, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	post := Post{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch {
		case util.IsElement(n, "title"):
			post.Title = util.GetNodeText(n)
			return
		case util.IsElement(n, "header"):
			// TODO: parseHeader
			return
		case util.IsElement(n, "section") && util.GetNodeAttr(n, "data-field") == "subtitle":
			// TODO: parserSubtitle
			return
		case util.IsElement(n, "section") && util.GetNodeAttr(n, "data-field") == "body":
			// TODO: parseBody
		case util.IsElement(n, "footer"):
			parseFooter(n, &post)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return &post, nil
}
