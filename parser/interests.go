package parser

import (
	"io"

	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func walkLinks(n *html.Node, f func(string, string)) {
	for t := n.FirstChild; t != nil; t = t.NextSibling {
		if t.Type == html.ElementNode && t.Data == "a" {
			f(util.GetNodeAttr(t, "href"), util.GetNodeText(t))
		}
	}
}

func ParseInterestsPublications(dat io.Reader) ([]Publication, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	pubs := []Publication{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if util.IsListItem(n) {
			p := Publication{}
			walkLinks(n, func(href string, text string) {
				p.Url = href
				p.Name = text
			})
			pubs = append(pubs, p)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return pubs, nil
}

func ParseInterestsTags(dat io.Reader) ([]Tag, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	tags := []Tag{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if util.IsListItem(n) {
			t := Tag{}
			walkLinks(n, func(href string, text string) {
				t.Url = href
				t.Name = text
			})
			tags = append(tags, t)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return tags, nil
}

func ParseInterestsTopics(dat io.Reader) ([]Topic, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	topics := []Topic{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if util.IsListItem(n) {
			t := Topic{}
			walkLinks(n, func(href string, text string) {
				t.Url = href
				t.Name = text
			})
			topics = append(topics, t)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return topics, nil
}

func ParseInterestsWriters(dat io.Reader) ([]User, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	users := []User{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if util.IsListItem(n) {
			u := User{}
			walkLinks(n, func(href string, text string) {
				u.Url = href
				u.Name = text
				u.Username = util.ParseMediumUsername(href)
			})
			users = append(users, u)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return users, nil
}
