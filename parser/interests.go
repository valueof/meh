package parser

import (
	"io"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func walkLinks(n *util.Node, f func(string, string)) {
	for t := n.FirstChild; t != nil; t = t.NextSibling {
		if t.IsElement("a") {
			f(t.Attrs["href"], t.Text())
		}
	}
}

func ParseInterestsPublications(dat io.Reader) ([]schema.Publication, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	pubs := []schema.Publication{}

	node.WalkChildren(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		p := schema.Publication{}
		walkLinks(n, func(href string, text string) {
			p.Url = href
			p.Name = text
		})
		pubs = append(pubs, p)
	})

	return pubs, nil
}

func ParseInterestsTags(dat io.Reader) ([]schema.Tag, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	tags := []schema.Tag{}

	node.WalkChildren(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		t := schema.Tag{}
		walkLinks(n, func(href string, text string) {
			t.Url = href
			t.Name = text
		})
		tags = append(tags, t)
	})

	return tags, nil
}

func ParseInterestsTopics(dat io.Reader) ([]schema.Topic, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	topics := []schema.Topic{}

	node.WalkChildren(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		t := schema.Topic{}
		walkLinks(n, func(href string, text string) {
			t.Url = href
			t.Name = text
		})
		topics = append(topics, t)
	})

	return topics, nil
}

func ParseInterestsWriters(dat io.Reader) ([]schema.User, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	users := []schema.User{}

	node.WalkChildren(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		u := schema.User{}
		walkLinks(n, func(href string, text string) {
			u.Url = href
			u.Name = text
			u.Username = util.ParseMediumUsername(href)
		})
		users = append(users, u)
	})

	return users, nil
}
