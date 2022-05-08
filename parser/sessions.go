package parser

import (
	"io"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func ParseSessions(dat io.Reader) (sessions []schema.Session, err error) {
	sessions = []schema.Session{}
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return
	}

	var f func(n *util.Node)
	f = func(n *util.Node) {
		if n.IsElement("li") {
			s := schema.Session{}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.TextNode {
					continue
				}

				// Medium exports this in an annoying format which, most likely, was implemented by me.
				// <ul>
				//  <li>
				//   Created at: <string>
				//   <br>
				//   Last seen at: <string>
				//   <br>
				//   Last seen location: <string>
				//   <br>
				//   User agent: <string>
				//  </li>
				// </ul>

				switch {
				case strings.HasPrefix(c.Data, "Created at:"):
					s.CreatedAt = strings.TrimPrefix(c.Text(), "Created at: ")
				case strings.HasPrefix(c.Data, "Last seen at:"):
					s.LastSeenAt = strings.TrimPrefix(c.Text(), "Last seen at: ")
				case strings.HasPrefix(c.Data, "Last seen location:"):
					s.LastSeenLocation = strings.TrimPrefix(c.Text(), "Last seen location: ")
				case strings.HasPrefix(c.Data, "User agent:"):
					s.UserAgent = strings.TrimPrefix(c.Text(), "User agent: ")
				}
			}

			sessions = append(sessions, s)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return
}
