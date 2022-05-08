package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func parseList(dat io.Reader, fn func(*util.Node)) error {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return err
	}

	var f func(n *util.Node)
	f = func(n *util.Node) {
		if n.IsElement("a") {
			fn(n)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return nil
}

func ParsePublicationFollowing(dat io.Reader) (pubs []schema.Publication, err error) {
	pubs = []schema.Publication{}
	err = parseList(dat, func(a *util.Node) {
		pubs = append(pubs, schema.Publication{
			Url:  a.Attrs["href"],
			Name: a.Text(),
		})
	})
	return
}

func ParseTopicsFollowing(dat io.Reader) (topics []schema.Topic, err error) {
	topics = []schema.Topic{}
	err = parseList(dat, func(a *util.Node) {
		topics = append(topics, schema.Topic{
			Url:  a.Attrs["href"],
			Name: a.Text(),
		})
	})
	return
}

func ParseUsersFollowing(dat io.Reader) (users []schema.User, err error) {
	users = []schema.User{}
	err = parseList(dat, func(a *util.Node) {
		users = append(users, schema.User{
			Url:      a.Attrs["href"],
			Username: strings.TrimPrefix(a.Text(), "@"),
		})
	})
	return
}

// ParseUsersSuggested parses suggested Twitter friends
//
// Medium sourced these suggestions from your Twitter account. These are
// users you follow on Twitter that are also on Medium.
func ParseUsersSuggested(dat io.Reader) (users []schema.User, err error) {
	users = []schema.User{}
	err = parseList(dat, func(a *util.Node) {
		users = append(users, schema.User{
			Url:      a.Attrs["href"],
			Username: strings.TrimPrefix(a.Text(), "@"),
		})
	})
	return
}
