package parser

import (
	"io"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func ParseHighlights(dat io.Reader) ([]schema.Highlight, error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	highlights := []schema.Highlight{}

	doc.WalkChildren(func(n *util.Node) {
		if !n.IsElement("li") || !n.HasClass("h-entry") {
			return
		}

		h := schema.Highlight{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.IsElement("time") && c.HasClass("dt-published") {
				h.CreatedAt = c.Text()
				break
			}
		}

		// Node.ParseGrafs ignores non-graf elements so we don't need to do any
		// additional parsing or stripping here.
		h.Body = n.ParseGrafs()
		highlights = append(highlights, h)
	})

	return highlights, nil
}
