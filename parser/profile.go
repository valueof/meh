package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
	"golang.org/x/net/html"
)

func ParseBio(dat io.Reader) (bio string, err error) {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return "", err
	}

	doc.WalkChildren(func(n *util.Node) {
		if n.IsElement("section") && n.HasClass("h-card") {
			bio = n.Text()
		}
	})

	return
}

func parseAccountInfo(ul *util.Node, p *schema.Profile) {
	for li := ul.FirstChildElement("li"); li != nil; li = li.NextSiblingElement("li") {
		b := li.FirstChildElement("b")
		switch b.Text() {
		case "Profile:":
			li.WalkChildren(func(c *util.Node) {
				if c.IsElement("a") && c.HasClass("u-url") {
					p.User.Url = c.Attrs["href"]
					p.User.Username = strings.TrimPrefix(c.Text(), "@")
				}
			})
		case "Display name:":
			p.User.Name = strings.TrimPrefix(li.Text(), "Display name: ")
		case "Email address:":
			p.Email = strings.TrimPrefix(li.Text(), "Email address: ")
		case "Previous email address:":
			p.PastEmails = []string{
				strings.TrimPrefix(li.Text(), "Previous email address: "),
			}
		case "Medium user ID:":
			p.User.Id = strings.TrimPrefix(li.Text(), "Medium user ID: ")
		case "Created at:":
			p.User.CreatedAt = strings.TrimPrefix(li.Text(), "Created at: ")
		}
	}
}

func parseSocialAccounts(ul *util.Node, p *schema.Profile) {
	twitter := schema.SocialAccount{}
	google := schema.SocialAccount{}

	for li := ul.FirstChildElement("li"); li != nil; li = li.NextSiblingElement("li") {
		b := li.FirstChildElement("b")
		switch b.Text() {
		case "Twitter:":
			li.WalkChildren(func(c *util.Node) {
				if c.IsElement("a") {
					twitter.Url = c.Attrs["href"]
					twitter.Name = c.Text()
				}
			})
		case "Twitter account ID:":
			twitter.Id = strings.TrimPrefix(li.Text(), "Twitter account ID: ")
		case "Google email:":
			google.Email = strings.TrimPrefix(li.Text(), "Google email: ")
		case "Google display name:":
			google.Name = strings.TrimPrefix(li.Text(), "Google display name: ")
		case "Google account ID:":
			google.Id = strings.TrimPrefix(li.Text(), "Google account ID: ")
		}

	}

	p.SocialAccounts = map[string]schema.SocialAccount{
		"twitter": twitter,
		"google":  google,
	}
}

func ParseUserProfile(dat io.Reader, profile *schema.Profile) error {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return err
	}

	doc.WalkChildren(func(n *util.Node) {
		if !n.IsElement("section") || !n.HasClass("h-card") {
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			switch {
			case c.IsElement("h3") && c.HasClass("p-name"):
				profile.User.Name = c.Text()
			case c.IsElement("img") && c.HasClass("u-photo"):
				img := schema.Image{
					Source: c.Attrs["src"],
				}
				profile.User.ProfilePic = &img
			case c.IsElement("h4") && c.Text() == "Account info":
				parseAccountInfo(c.NextSiblingElement("ul"), profile)
			case c.IsElement("h4") && c.Text() == "Connected accounts":
				parseSocialAccounts(c.NextSiblingElement("ul"), profile)
			}
		}
	})

	return nil
}

func parsePubs(ul *util.Node) []schema.Publication {
	pubs := []schema.Publication{}
	for li := ul.FirstChildElement("li"); li != nil; li = li.NextSiblingElement("li") {
		a := li.FirstChildElement("a")
		pubs = append(pubs, schema.Publication{
			Url:  a.Attrs["href"],
			Name: a.Text(),
		})
	}
	return pubs
}

func ParsePublications(dat io.Reader, profile *schema.Profile) error {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return err
	}

	doc.WalkChildren(func(n *util.Node) {
		if !n.IsElement("section") {
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			switch {
			case c.IsElement("h4") && c.Text() == "Editor":
				profile.Editor = parsePubs(c.NextSiblingElement("ul"))
			case c.IsElement("h4") && c.Text() == "Writer":
				profile.Writer = parsePubs(c.NextSiblingElement("ul"))
			}
		}
	})

	return nil
}

func ParseMemberships(dat io.Reader, profile *schema.Profile) error {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return err
	}

	profile.Memberships = []schema.Membership{}
	doc.WalkChildren(func(n *util.Node) {
		if !n.IsElement("ul") {
			return
		}

		for li := n.FirstChildElement("li"); li != nil; li = li.NextSiblingElement("li") {
			m := schema.Membership{}

			for t := li.FirstChild; t != nil; t = t.NextSibling {
				if t.Type != html.TextNode {
					continue
				}

				switch {
				case strings.HasPrefix(t.Data, "Membership ID:"):
					m.Id = strings.TrimPrefix(t.Data, "Membership ID: ")
				case strings.HasPrefix(t.Data, "Started at:"):
					m.StartedAt = strings.TrimPrefix(t.Data, "Started at: ")
				case strings.HasPrefix(t.Data, "Ended at:"):
					m.EndedAt = strings.TrimPrefix(t.Data, "Ended at: ")
				case strings.HasPrefix(t.Data, "Amount:"):
					m.Amount, _ = strconv.ParseFloat(strings.TrimPrefix(t.Data, "Amount: $"), 64)
				case strings.HasPrefix(t.Data, "Type:"):
					m.EndedAt = strings.TrimPrefix(t.Data, "Type: ")
				}
			}

			profile.Memberships = append(profile.Memberships, m)
		}
	})

	return nil
}

func ParseMembershipCharges(dat io.Reader, profile *schema.Profile) error {
	doc, err := util.NewNodeFromHTML(dat)
	if err != nil {
		return err
	}

	if profile.MembershipCharges == nil {
		profile.MembershipCharges = []schema.MembershipCharge{}
	}

	doc.WalkChildren(func(n *util.Node) {
		if !n.IsElement("ul") {
			return
		}

		for li := n.FirstChildElement("li"); li != nil; li = li.NextSiblingElement("li") {
			c := schema.MembershipCharge{}

			for t := li.FirstChild; t != nil; t = t.NextSibling {
				if t.Type != html.TextNode {
					continue
				}

				switch {
				case strings.HasPrefix(t.Data, "Created at:"):
					c.CreatedAt = strings.TrimPrefix(t.Data, "Created at: ")
				case strings.HasPrefix(t.Data, "Amount:"):
					c.Amount, _ = strconv.ParseFloat(strings.TrimPrefix(t.Data, "Amount: $"), 64)
				}
			}

			profile.MembershipCharges = append(profile.MembershipCharges, c)
		}
	})

	return nil
}
