package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func isUser(n *html.Node) bool {
	if n.Type != html.ElementNode || n.Data != "a" {
		return false
	}

	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == "h-cite" {
			return true
		}
	}

	return false
}

// ParseBlocked parses the blocked-users HTML file
func ParseBlocked(dat io.Reader) ([]MediumUser, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	users := []MediumUser{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if isUser(n) {
			for t := n.FirstChild; t != nil; t = t.NextSibling {
				if t.Type == html.TextNode {
					users = append(users, MediumUser{
						Username: strings.TrimPrefix(t.Data, "@"),
					})
					break
				}
			}

			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return users, nil
}
