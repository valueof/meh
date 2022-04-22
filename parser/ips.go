package parser

import (
	"io"
	"net"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func ParseIps(dat io.Reader) ([]schema.IP, error) {
	node, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return nil, err
	}

	ips := []schema.IP{}

	node.WalkChildren(func(n *util.Node) {
		if n.IsElement("li") == false {
			return
		}

		ip := schema.IP{}
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
	})

	return ips, nil
}
