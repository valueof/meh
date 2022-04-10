package parser

import (
	"io"
	"net"
	"strings"

	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func ParseIps(dat io.Reader) ([]IP, error) {
	doc, err := html.Parse(dat)
	if err != nil {
		return nil, err
	}

	ips := []IP{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if util.IsListItem(n) {
			ip := IP{}
			for t := n.FirstChild; t != nil; t = t.NextSibling {
				if t.Type != html.TextNode {
					continue
				}

				switch {
				case strings.HasPrefix(t.Data, "IP:"):
					s := net.ParseIP(strings.TrimSpace(strings.TrimPrefix(t.Data, "IP:")))
					if s != nil {
						ip.Address = s.String()
					}
				case strings.HasPrefix(t.Data, "Created at:"):
					ip.CreatedAt = strings.TrimSpace(strings.TrimPrefix(t.Data, "Created at:"))
				}
			}
			if ip.Address != "" {
				ips = append(ips, ip)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return ips, nil
}
