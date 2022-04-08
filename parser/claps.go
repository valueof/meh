package parser

import (
	"io"
	"regexp"
	"strconv"

	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

var numRe regexp.Regexp = *regexp.MustCompile(`^.*\+(\d+) —`)

func ParseClaps(dat io.Reader) ([]Clap, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	claps := []Clap{}

	isE := func(n *html.Node) bool {
		return n.Type == html.ElementNode
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if isE(n) && n.Data == "li" {
			clap := Clap{}
			post := Post{}
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

				if isE(c) && util.GetNodeAttr(c, "class") == "h-cite u-like-of" {
					post.Url = util.GetNodeAttr(c, "href")
					post.Id = util.ParseMediumId(post.Url)
					post.Title = util.GetNodeText(c)
				}

				if isE(c) && util.GetNodeAttr(c, "class") == "dt-published" {
					post.PublishedAt = util.GetNodeText(c)
				}
			}
			clap.Post = post
			claps = append(claps, clap)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return claps, nil
}
