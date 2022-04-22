package parser

import (
	"io"
	"regexp"
	"strconv"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

var numRe regexp.Regexp = *regexp.MustCompile(`^.*\+(\d+) —`)

func ParseClaps(dat io.Reader) ([]schema.Clap, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	claps := []schema.Clap{}

	node.WalkChildren(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		clap := schema.Clap{}
		post := schema.Post{}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			// Parse out the first part of item:
			// 	+1 &mdash;
			// Number of claps can go from 1 to 50
			// &mdash; gets converted to — by html.Parse
			if c.Type == html.TextNode {
				m := numRe.FindStringSubmatch(c.Data)
				if len(m) > 1 {
					num, _ := strconv.Atoi(m[1])
					clap.Amount = num
				}
			}

			if c.IsElement("a") && c.Attrs["class"] == "h-cite u-like-of" {
				post.Url = c.Attrs["href"]
				post.Id = util.ParseMediumId(post.Url)
				post.Title = c.Text()
			}

			if c.IsElement("time") && c.Attrs["class"] == "dt-published" {
				post.PublishedAt = c.Text()
			}
		}
		clap.Post = post
		claps = append(claps, clap)
	})

	return claps, nil
}
