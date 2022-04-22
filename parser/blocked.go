package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

// ParseBlocked parses the blocked-users HTML file
func ParseBlocked(dat io.Reader) ([]schema.User, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	users := []schema.User{}
	node.WalkChildren(func(n *util.Node) {
		if n.Type == html.ElementNode && n.Data == "a" && n.Attrs["class"] == "h-cite" {
			users = append(users, schema.User{
				Username: strings.TrimPrefix(n.Text(), "@"),
			})
		}
	})

	return users, nil
}
