package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

// ParseBlocked parses the blocked-users HTML file
func ParseBlocked(dat io.Reader) ([]User, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	users := []User{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" && util.GetNodeAttr(n, "class") == "h-cite" {
			users = append(users, User{
				Username: strings.TrimPrefix(util.GetNodeText(n), "@"),
			})
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return users, nil
}
