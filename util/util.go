package util

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// ParseMediumId Parses post ID out of a Medium URL. Links to all Medium posts
// end with a unique value that represents its ID:
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

func GetNodeText(n *html.Node) string {
	s := []string{}

	for t := n.FirstChild; t != nil; t = t.NextSibling {
		if t.Type == html.TextNode {
			s = append(s, t.Data)
		}
	}

	out := strings.Join(s, "")
	re := regexp.MustCompile(`\s+`)
	out = re.ReplaceAllString(out, " ")
	return strings.TrimSpace(out)
}

func GetNodeAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}
