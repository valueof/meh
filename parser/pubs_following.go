package parser

import (
	"io"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func ParsePubsFollowing(dat io.Reader) (pubs []schema.Publication, err error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	pubs = []schema.Publication{}
	var f func(n *util.Node)
	f = func(n *util.Node) {
		if n.IsElement("a") {
			pubs = append(pubs, schema.Publication{
				Url:  n.Attrs["href"],
				Name: n.Text(),
			})
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return
}
