package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func parseBody(n *html.Node, post *Post) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch {
		case util.IsElement(c, "section"):
			post.Content = append(post.Content, Section{
				Name: util.GetNodeAttr(c, "name"),
				Body: parseInnerSections(c),
			})
		}
	}
}

func parseGrafs(body *html.Node) []Graf {
	grafs := []Graf{}

	for g := body.FirstChild; g != nil; g = g.NextSibling {
		if util.HasClass(g, "graf") == false {
			continue
		}

		graf := Graf{
			Name: util.GetNodeAttr(g, "name"),
			Type: GrafType(g.Data),
			Text: util.GetNodeAllText(g),
		}

		grafs = append(grafs, graf)
	}

	return grafs
}

func parseInnerSections(body *html.Node) []InnerSection {
	sections := []InnerSection{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if util.HasClass(n, "section-inner") {
			sub := InnerSection{
				Classes: []string{},
				Body:    parseGrafs(n),
			}

			for _, class := range strings.Split(util.GetNodeAttr(n, "class"), " ") {
				if class != "section-inner" {
					sub.Classes = append(sub.Classes, class)
				}
			}

			sections = append(sections, sub)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(body)
	return sections
}

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
		case util.IsElement(n, "section") && util.GetNodeAttr(n, "data-field") == "body":
			parseBody(n, &post)
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
