package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/util"
)

func parseBody(n *util.Node, post *Post) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch {
		case c.IsElement("section"):
			post.Content = append(post.Content, Section{
				Name: c.Attrs["name"],
				Body: parseInnerSections(c),
			})
		}
	}
}

func parseGrafs(body *util.Node) []Graf {
	grafs := []Graf{}

	for g := body.FirstChild; g != nil; g = g.NextSibling {
		if g.HasClass("graf") == false {
			continue
		}

		graf := Graf{Name: g.Attrs["name"]}

		// TODO(anton):
		// 	graf--pullquote
		//	graf--pre
		//	graf--sectionCaption
		//	markups

		switch {
		case g.HasClass("graf--h1"):
		case g.HasClass("graf--h2"):
		case g.HasClass("graf--h3"):
		case g.HasClass("graf--h4"):
		case g.HasClass("graf--blockquote"):
		case g.HasClass("graf--p"):
			graf.Type = GrafType(g.Data)
			graf.Text = g.Text()
		case g.HasClass("graf--figure"):
			graf.Type = IMG
			graf.Image = extractImage(g)
		case g.HasClass("graf--mixtapeEmbed"):
			graf.Type = EMBED
			graf.Text = g.Text()
			// TODO(anton): Better support for mixtapes
		}

		grafs = append(grafs, graf)
	}

	return grafs
}

func extractImage(n *util.Node) (img *Image) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.IsElement("img") == false {
			continue
		}

		return &Image{
			Name:   c.Attrs["data-image-id"],
			Width:  c.Attrs["data-width"],
			Height: c.Attrs["data-height"],
			Source: c.Attrs["src"],
		}
	}

	return nil
}

func parseInnerSections(body *util.Node) []InnerSection {
	sections := []InnerSection{}

	var f func(*util.Node)
	f = func(n *util.Node) {
		if n.HasClass("section-inner") {
			sub := InnerSection{
				Classes: []string{},
				Body:    parseGrafs(n),
			}

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

func parseFooter(n *util.Node, post *Post) {
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

func ParsePost(dat io.Reader) (*Post, error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	post := Post{}

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
