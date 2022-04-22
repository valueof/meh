package util

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/valueof/meh/schema"
	"golang.org/x/net/html"
)

var SPACE_RE *regexp.Regexp = regexp.MustCompile(`\s+`)

// ParseMediumId Parses post ID out of a Medium URL. Links to all Medium posts end with a unique
// value that represents its ID:
// 	https://medium.com/p/my-slug-5940ded906e7 -> 5940ded906e7
func ParseMediumId(s string) string {
	url, err := url.Parse(s)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile("-([a-z0-9]+)$")
	m := re.FindStringSubmatch(url.Path)
	if len(m) >= 2 {
		return m[1]
	}

	return ""
}

// ParseMediumUsername Parses username out of a Medium URL. For now it only supports
// medium.com/@username and username.medium.com.
//
// Caveat: sometimes username.medium.com is not username at all but we will ignore this
// fact for now.
func ParseMediumUsername(s string) string {
	url, err := url.Parse(s)
	if err != nil {
		return ""
	}

	p := strings.Split(url.Path, "/")
	if len(p) > 1 && strings.HasPrefix(p[1], "@") {
		return strings.TrimPrefix(p[1], "@")
	}

	h := strings.Split(url.Host, ".")
	if len(h) > 2 {
		return h[0]
	}

	return ""
}

type Node struct {
	*html.Node
	FirstChild  *Node
	NextSibling *Node
	Attrs       map[string]string
}

func NewNodeFromHTML(dat io.Reader) (*Node, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	return NewNode(doc), nil
}

func NewNode(n *html.Node) *Node {
	node := Node{
		Node:        n,
		FirstChild:  nil,
		NextSibling: nil,
		Attrs:       map[string]string{},
	}

	if n.FirstChild != nil {
		node.FirstChild = NewNode(n.FirstChild)
	}

	if n.NextSibling != nil {
		node.NextSibling = NewNode(n.NextSibling)
	}

	for _, a := range n.Attr {
		node.Attrs[a.Key] = a.Val
	}

	return &node
}

func (n *Node) WalkChildren(cb func(*Node)) {
	var f func(*Node)
	f = func(n *Node) {
		cb(n)

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f(c)
	}
}

func (n *Node) Text() (s string) {
	s = ""
	n.WalkChildren(func(t *Node) {
		if t.Type == html.TextNode {
			s += trimUnnecessarySpace(t.Data)
		}
	})
	s = strings.TrimSpace(s)
	return
}

func (n *Node) Markup() (markup []schema.Markup) {
	s := ""
	markup = []schema.Markup{}

	n.WalkChildren(func(t *Node) {
		fn := func(tp schema.MarkupType) schema.Markup {
			start := len(s)
			end := start + len(t.Text())
			return schema.Markup{Type: tp, Start: start, End: end}
		}

		if t.Type == html.TextNode {
			ns := trimUnnecessarySpace(t.Data)
			if s == "" { // Beginning of this node
				ns = strings.TrimPrefix(ns, " ")
			}
			s += ns
			return
		}

		if t.Type != html.ElementNode {
			return
		}

		switch t.Data {
		case "em":
			markup = append(markup, fn(schema.EM))
		case "strong":
			markup = append(markup, fn(schema.STRONG))
		case "br":
			markup = append(markup, fn(schema.BR))
		case "a":
			a := fn(schema.A)
			a.Href = t.Attrs["href"]
			markup = append(markup, a)

		default:
			fmt.Printf("Unknown markup: %s; %s\n", t.Data, t.Text())
		}
	})

	return
}

func (n *Node) TextPreformatted() string {
	s := []string{}

	n.WalkChildren(func(t *Node) {
		if t.Type == html.TextNode {
			s = append(s, t.Data)
		} else if t.IsElement("br") {
			s = append(s, "\n")
		}
	})

	return strings.Join(s, "")
}

func (n *Node) IsElement(name string) bool {
	return n.Type == html.ElementNode && n.Data == name
}

func (n *Node) HasClass(name string) bool {
	if n.Type != html.ElementNode {
		return false
	}

	classes := strings.Split(n.Attrs["class"], " ")
	for _, c := range classes {
		if c == name {
			return true
		}
	}

	return false
}

func trimUnnecessarySpace(s string) string {
	ns := SPACE_RE.ReplaceAllString(s, " ")
	if strings.TrimSpace(ns) != "" {
		return ns
	}
	return ""
}
