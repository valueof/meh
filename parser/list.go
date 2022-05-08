package parser

import (
	"io"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func ParseList(dat io.Reader) (lists *schema.List, err error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	lists = &schema.List{
		Name:    "",
		Summary: "",
		Posts:   []schema.Post{},
	}

	var f func(*util.Node)
	f = func(n *util.Node) {
		if n.IsElement("h1") && n.HasClass("p-name") {
			lists.Name = n.Text()
			return
		}

		if n.IsElement("h2") && n.HasClass("p-summary") {
			lists.Summary = n.Text()
			return
		}

		if n.IsElement("li") && n.Attrs["data-field"] == "post" {
			for a := n.FirstChild; a != nil; a = a.NextSibling {
				if !a.IsElement("a") {
					continue
				}

				lists.Posts = append(lists.Posts, schema.Post{
					Id:    util.ParseMediumId(a.Attrs["href"]),
					Url:   a.Attrs["href"],
					Title: a.Text(),
				})
			}

			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return
}
