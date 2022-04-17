package parser

import (
	"io"

	"github.com/valueof/meh/util"
)

func walkLinks(n *util.Node, f func(string, string)) {
	for t := n.FirstChild; t != nil; t = t.NextSibling {
		if t.IsElement("a") {
			f(t.Attrs["href"], t.Text())
		}
	}
}

func ParseInterestsPublications(dat io.Reader) ([]Publication, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	pubs := []Publication{}

	node.Walk(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		p := Publication{}
		walkLinks(n, func(href string, text string) {
			p.Url = href
			p.Name = text
		})
		pubs = append(pubs, p)
	})

	return pubs, nil
}

func ParseInterestsTags(dat io.Reader) ([]Tag, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	tags := []Tag{}

	node.Walk(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		t := Tag{}
		walkLinks(n, func(href string, text string) {
			t.Url = href
			t.Name = text
		})
		tags = append(tags, t)
	})

	return tags, nil
}

func ParseInterestsTopics(dat io.Reader) ([]Topic, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	topics := []Topic{}

	node.Walk(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		t := Topic{}
		walkLinks(n, func(href string, text string) {
			t.Url = href
			t.Name = text
		})
		topics = append(topics, t)
	})

	return topics, nil
}

func ParseInterestsWriters(dat io.Reader) ([]User, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	users := []User{}

	node.Walk(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		u := User{}
		walkLinks(n, func(href string, text string) {
			u.Url = href
			u.Name = text
			u.Username = util.ParseMediumUsername(href)
		})
		users = append(users, u)
	})

	return users, nil
}
