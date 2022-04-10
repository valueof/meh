package util

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// ParseMediumId Parses post ID out of a Medium URL. Links to all Medium
// posts end with a unique value that represents its ID:
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

// ParseMediumUsername Parses username out of a Medium URL. For now it
// only supports medium.com/@username and username.medium.com.
//
// Caveat: sometimes username.medium.com is not username at all
// but we will ignore this fact for now. If you think this is confusing
// ask someone from Medium about difference between publications, collections,
// and catalogs and watch them weep.
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

func IsListItem(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "li"
}
