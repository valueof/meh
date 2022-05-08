package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func ParsePublicationFollowing(dat io.Reader) (pubs []schema.Publication, err error) {
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

func ParseTopicsFollowing(dat io.Reader) (pubs []schema.Topic, err error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	pubs = []schema.Topic{}
	var f func(n *util.Node)
	f = func(n *util.Node) {
		if n.IsElement("a") {
			pubs = append(pubs, schema.Topic{
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

func ParseUsersFollowing(dat io.Reader) (users []schema.User, err error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	users = []schema.User{}

	var f func(n *util.Node)
	f = func(n *util.Node) {
		if n.IsElement("a") {
			users = append(users, schema.User{
				Url:      n.Attrs["href"],
				Username: strings.TrimPrefix(n.Text(), "@"),
			})
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return
}
