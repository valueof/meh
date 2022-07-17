package parser

import (
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/valueof/meh/formatters"
	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

type Parser struct {
	logger    log.Logger
	root      string
	formatter formatters.Formatter
}

func NewParser(root string, logger log.Logger, f formatters.Formatter) *Parser {
	return &Parser{
		logger:    logger,
		root:      root,
		formatter: f,
	}
}

func (p *Parser) walk(d fs.FileInfo, fn func(string, io.Reader)) error {
	dir := filepath.Join(p.root, d.Name())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		p.logger.Printf("can't read %s, skipping", path.Base(dir))
		return err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") == false {
			p.logger.Printf("%s is not an html file, skipping", f.Name())
			continue
		}

		dat, err := os.Open(path.Join(dir, f.Name()))
		if err != nil {
			p.logger.Printf("can't read %s, skipping", f.Name())
			continue
		}
		defer dat.Close()

		fn(f.Name(), dat)
	}

	return nil
}

func (p *Parser) Parse() error {
	dirs, err := ioutil.ReadDir(p.root)
	if err != nil {
		return err
	}

	for _, d := range dirs {
		if d.IsDir() == false {
			p.logger.Printf("%s is not a directory, skipping", d.Name())
			continue
		}

		switch d.Name() {
		case "blocks":
			users := []schema.User{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseBlocked(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				users = append(users, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("blocks", schema.BlockedUsers{
				Meta:  "Blocked users",
				Users: users,
			})
		case "bookmarks":
			posts := []schema.Post{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseBookmarks(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				posts = append(posts, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("bookmarks", schema.Bookmarks{
				Meta:  "Bookmarked posts",
				Posts: posts,
			})
		case "claps":
			claps := []schema.Clap{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseClaps(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				claps = append(claps, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("claps", schema.Claps{
				Meta:  "Posts you've clapped for",
				Claps: claps,
			})
		case "interests":
			interests := schema.Interests{
				Meta: "Topics you're interested in",
			}

			err = p.walk(d, func(name string, dat io.Reader) {
				switch name {
				case "publications.html":
					pubs, err := ParseInterestsPublications(dat)
					if err != nil {
						p.logger.Printf("error parsing %s, skipping", name)
						return
					}
					p.logger.Printf("parsed %s", name)
					interests.Publications = pubs
				case "tags.html":
					tags, err := ParseInterestsTags(dat)
					if err != nil {
						p.logger.Printf("error parsing %s, skipping", name)
						return
					}
					p.logger.Printf("parsed %s", name)
					interests.Tags = tags
				case "topics.html":
					topics, err := ParseInterestsTopics(dat)
					if err != nil {
						p.logger.Printf("error parsing %s, skipping", name)
						return
					}
					p.logger.Printf("parsed %s", name)
					interests.Topics = topics
				case "writers.html":
					writers, err := ParseInterestsWriters(dat)
					if err != nil {
						p.logger.Printf("error parsing %s, skipping", name)
						return
					}
					p.logger.Printf("parsed %s", name)
					interests.Writers = writers
				default:
					p.logger.Printf("Unknown interests file %s, skipping", name)
				}
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("interests", interests)
		case "ips":
			ips := []schema.IP{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseIps(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				ips = append(ips, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("ips", schema.IPs{
				Meta: "Your IP history (note: Medium deletes IP history after 30 days)",
				IPs:  ips,
			})
		case "posts":
			posts := map[string]schema.Post{}
			err = p.walk(d, func(name string, dat io.Reader) {
				post, err := ParsePost(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				posts[strings.TrimSuffix(name, ".html")] = *post
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			for name, post := range posts {
				p.formatter.WriteFile(filepath.Join("posts", name), post)
			}
		case "lists":
			lists := []schema.List{}
			err = p.walk(d, func(name string, dat io.Reader) {
				list, err := ParseList(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				lists = append(lists, *list)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("lists", schema.Lists{
				Meta:  "Lists you've created",
				Lists: lists,
			})
		case "pubs-following":
			pubs := []schema.Publication{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParsePublicationFollowing(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				pubs = append(pubs, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("following/publications", schema.Publications{
				Meta:         "Publications you follow",
				Publications: pubs,
			})
		case "topics-following":
			topics := []schema.Topic{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseTopicsFollowing(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				topics = append(topics, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("following/topics", schema.Topics{
				Meta:   "Topics you follow",
				Topics: topics,
			})
		case "users-following":
			users := []schema.User{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseUsersFollowing(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				users = append(users, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile(filepath.Join("following", "users"), schema.Users{
				Meta:  "Users you follow",
				Users: users,
			})
		case "twitter":
			users := []schema.User{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseUsersSuggested(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				users = append(users, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile(filepath.Join("following", "suggested"), schema.Users{
				Meta:  "Your Twitter friends who are also on Medium",
				Users: users,
			})
		case "sessions":
			sessions := []schema.Session{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseSessions(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				sessions = append(sessions, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("sessions", schema.Sessions{
				Meta:     "Your active and inactive sessions across devices",
				Sessions: sessions,
			})
		case "highlights":
			highlights := []schema.Highlight{}
			err = p.walk(d, func(name string, dat io.Reader) {
				part, err := ParseHighlights(dat)
				if err != nil {
					p.logger.Printf("error parsing %s, skipping", name)
					return
				}
				p.logger.Printf("parsed %s", name)
				highlights = append(highlights, part...)
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("highlights", schema.Highlights{
				Meta:       "Your highlights",
				Highlights: highlights,
			})
		case "profile":
			profile := schema.Profile{
				Meta: "Your user profile",
			}
			profile.User = &schema.User{}

			err = p.walk(d, func(name string, dat io.Reader) {
				switch {
				case name == "about.html":
					bio, err := ParseBio(dat)
					if err != nil {
						p.logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					p.logger.Printf("parsed %s", name)
					profile.User.Bio = bio
				case name == "profile.html":
					err = ParseUserProfile(dat, &profile)
					if err != nil {
						p.logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					p.logger.Printf("parsed %s", name)
				case name == "publications.html":
					err = ParsePublications(dat, &profile)
					if err != nil {
						p.logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					p.logger.Printf("parsed %s", name)
				case name == "memberships.html":
					err = ParseMemberships(dat, &profile)
					if err != nil {
						p.logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					p.logger.Printf("parsed %s", name)
				case strings.HasPrefix(name, "charges-") && strings.HasSuffix(name, ".html"):
					err = ParseMembershipCharges(dat, &profile)
					if err != nil {
						p.logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					p.logger.Printf("parsed %s", name)
				default:
					p.logger.Printf("skipped profile/%s: not supported", name)
				}
			})

			if err != nil {
				p.logger.Printf("error parsing %s: %v", d.Name(), err)
				continue
			}

			p.formatter.WriteFile("profile", profile)
		default:
			p.logger.Printf("%s isn't supported, skipping", d.Name())
		}
	}

	return nil
}

func (p *Parser) FetchImages(dest string) {
	dir := filepath.Join(dest, "images")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		p.logger.Printf("couldn't create %s, images will not be downloaded", dir)
		return
	}

	var wg sync.WaitGroup
	for _, img := range util.GetQueuedImages() {
		wg.Add(1)
		img := img
		go func() {
			defer wg.Done()
			err := util.DownloadImage(img, filepath.Join(dir, img))
			if err != nil {
				p.logger.Printf("error downloading image %s. err: %v", img, err)
			}
			p.logger.Printf("downloaded %s", img)
		}()
	}
	wg.Wait()
}
