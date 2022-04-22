package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func parseBody(n *util.Node, post *schema.Post) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch {
		case c.IsElement("section"):
			post.Content = append(post.Content, schema.Section{
				Name: c.Attrs["name"],
				Body: parseInnerSections(c),
			})
		}
	}
}

func parseGrafs(body *util.Node) []schema.Graf {
	grafs := []schema.Graf{}

	for g := body.FirstChild; g != nil; g = g.NextSibling {
		if g.HasClass("graf") == false {
			continue
		}

		graf := schema.Graf{
			Name:    g.Attrs["name"],
			Markups: []schema.Markup{},
		}

		switch {
		case g.HasClass("graf--h1"):
			graf.Type = schema.H1
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--h2"):
			graf.Type = schema.H2
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--h3"):
			graf.Type = schema.H3
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--h4"):
			graf.Type = schema.H4
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--blockquote"):
			fallthrough
		case g.HasClass("graf--pullquote"):
			graf.Type = schema.BLOCKQUOTE
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--p"):
			graf.Type = schema.P
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--figure"):
			graf.Type = schema.IMG
			graf.Image = extractImage(g)
		case g.HasClass("graf--mixtapeEmbed"):
			graf.Type = schema.EMBED
			graf.Text = g.Text()
			// TODO(anton): Better support for mixtapes
		case g.HasClass("graf--pre"):
			graf.Type = schema.PRE
			graf.Text = g.TextPreformatted()
		case g.HasClass("graf--empty"):
			// Ignore empty grafs
		default:
			fmt.Printf("Unknown graf type: %s\n", g.Attrs["class"])
		}

		if graf.Type != "" {
			grafs = append(grafs, graf)
		}
	}

	return grafs
}

func extractImage(n *util.Node) (img *schema.Image) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.IsElement("img") == false {
			continue
		}

		return &schema.Image{
			Name:   c.Attrs["data-image-id"],
			Width:  c.Attrs["data-width"],
			Height: c.Attrs["data-height"],
			Source: c.Attrs["src"],
		}
	}

	return nil
}

func parseInnerSections(body *util.Node) []schema.InnerSection {
	sections := []schema.InnerSection{}

	var f func(*util.Node)
	f = func(n *util.Node) {
		if n.HasClass("section-inner") {
			grafs := parseGrafs(n)
			if len(grafs) == 0 {
				return
			}

			sub := schema.InnerSection{Body: grafs, Classes: []string{}}
			for _, class := range strings.Split(n.Attrs["class"], " ") {
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

func parseFooter(n *util.Node, post *schema.Post) {
	switch {
	case n.IsElement("time") && n.HasClass("dt-published"):
		post.PublishedAt = n.Attrs["datetime"]
	case n.IsElement("a") && n.HasClass("p-canonical"):
		post.Url = n.Attrs["href"]
		post.Id = util.ParseMediumId(post.Url)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseFooter(c, post)
	}
}

func ParsePost(dat io.Reader) (*schema.Post, error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	post := schema.Post{}

	var f func(*util.Node)
	f = func(n *util.Node) {
		switch {
		case n.IsElement("title"):
			post.Title = n.Text()
			return
		case n.IsElement("section") && n.Attrs["data-field"] == "body":
			parseBody(n, &post)
		case n.IsElement("footer"):
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
