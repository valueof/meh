/*
Package util implements functions useful when parsing and traversing HTML
generated by Medium's export tool.
*/
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

// ParseMediumId Parses post ID out of a Medium URL. Links to all Medium posts
// end with a unique value that represents its ID:
// 	https://medium.com/p/my-slug-5940ded906e7 -> 5940ded906e7
// 	https://medium.com/p/5940ded906e7         -> 5940ded906e7
func ParseMediumId(s string) string {
	url, err := url.Parse(s)
	if err != nil {
		return ""
	}

	re1 := regexp.MustCompile("-([a-z0-9]+)$")
	re2 := regexp.MustCompile(`\/p\/([a-z0-9]+)$`)

	m := re1.FindStringSubmatch(url.Path)
	if len(m) >= 2 {
		return m[1]
	}

	m = re2.FindStringSubmatch(url.Path)
	if len(m) >= 2 {
		return m[1]
	}

	return ""
}

// ParseMediumUsername Parses username out of a Medium URL. For now it only
// supports medium.com/@username and username.medium.com.
//
// Caveat: sometimes username.medium.com is not username at all but we will
// ignore this fact for now.
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

// A version of html.Node with easier access to attributes via Attrs
type Node struct {
	*html.Node
	FirstChild  *Node
	NextSibling *Node
	Attrs       map[string]string
}

// NewNodeFromHTML returns the parse tree for the HTML from the given Reader.
//
// Under the hood it uses html.Parse to parse HTML but wraps the result into
// the util.Node.
func NewNodeFromHTML(dat io.Reader) (*Node, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	return NewNode(doc), nil
}

// NewNode wraps html.Node into util.Node
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

// WalkChildren does a depth-first walk through all children of a given Node
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

// Text returns text content of a given Node, stripping away all HTML elements.
// It also collapses unnecessary space into a single space character and trims
// space at the beginning and end.
//
// Examples:
// 	In:  <p>This is a message from the <strong>log</strong></p>
//  Out: This is a message from the log
//
//  In:  <p>
//         This is a message from the
//           <strong>
//             log
//           </strong>
//       </p>
//	Out: This is a message from the log
func (n *Node) Text() (s string) {
	s = ""

	if n.Type == html.TextNode {
		s = strings.TrimSpace(collapseSpace(n.Data))
	} else {
		n.WalkChildren(func(t *Node) {
			if t.Type == html.TextNode {
				s += collapseSpace(t.Data)
			}
		})
		s = strings.TrimSpace(s)
	}
	return
}

// TextPreformatted is like Text except it preserves all spaces
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

// Markup returns a stacked slice of schema.Markup for the giving Node
// relative (and applicable to) the output of Text()
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
			ns := collapseSpace(t.Data)
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
		case "span":
			if t.HasClass("markup--highlight") {
				markup = append(markup, fn(schema.HIGHLIGHT))
			}
		default:
			fmt.Printf("Unknown markup: %s; %s\n", t.Data, t.Text())
		}
	})

	return
}

// Extract extracts image metadata from a given Node
func (n *Node) ExtractImage() (img *schema.Image) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.IsElement("img") == false {
			continue
		}

		return &schema.Image{
			Name:   c.Attrs["data-image-id"],
			Width:  c.Attrs["data-width"],
			Height: c.Attrs["data-height"],
			Source: c.Attrs["src"],
		}
	}

	return nil
}

// ParseGrafs parses a give Node and extracts all grafs, together with their markups.
func (n *Node) ParseGrafs() []schema.Graf {
	grafs := []schema.Graf{}

	for g := n.FirstChild; g != nil; g = g.NextSibling {
		if g.HasClass("graf") == false {
			continue
		}

		graf := schema.Graf{
			Name:    g.Attrs["name"],
			Markups: []schema.Markup{},
		}

		switch {
		case g.HasClass("graf--h1"):
			graf.Type = schema.H1
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--h2"):
			graf.Type = schema.H2
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--h3"):
			graf.Type = schema.H3
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--h4"):
			graf.Type = schema.H4
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--blockquote"):
			fallthrough
		case g.HasClass("graf--pullquote"):
			graf.Type = schema.BLOCKQUOTE
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--p"):
			graf.Type = schema.P
			graf.Text = g.Text()
			graf.Markups = g.Markup()
		case g.HasClass("graf--figure"):
			graf.Type = schema.IMG
			graf.Image = g.ExtractImage()
		case g.HasClass("graf--mixtapeEmbed"):
			graf.Type = schema.EMBED
			graf.Text = g.Text()
			// TODO(anton): Better support for mixtapes
		case g.HasClass("graf--pre"):
			graf.Type = schema.PRE
			graf.Text = g.TextPreformatted()
		case g.HasClass("graf--empty"):
			// Ignore empty grafs
		default:
			fmt.Printf("Unknown graf type: %s\n", g.Attrs["class"])
		}

		if graf.Type != "" {
			grafs = append(grafs, graf)
		}
	}

	return grafs
}

// IsElement returns true if the Node is html.ElementNode with a given tag name
func (n *Node) IsElement(name string) bool {
	return n.Type == html.ElementNode && n.Data == name
}

// IsElement returns true if the Node contains a given class
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

// Collapses spaces into one. If the input string contains only
// spaces returns an empty string. Can als reverse the Big Bang.
func collapseSpace(s string) string {
	ns := SPACE_RE.ReplaceAllString(s, " ")
	if strings.TrimSpace(ns) != "" {
		return ns
	}
	return ""
}
