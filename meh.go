package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/schema"
)

var input *string
var output *string
var verbose *bool

func walk(dir string, logger *log.Logger, fn func(string, io.Reader)) {
	dir = path.Join(*input, dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Fatalf("%s: %v\n", path.Base(dir), err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") == false {
			logger.Printf("skipped %s: not .html", f.Name())
			continue
		}

		dat, err := os.Open(path.Join(dir, f.Name()))
		defer dat.Close()
		if err != nil {
			logger.Fatalf("%s: %v", f.Name(), err)
			continue
		}

		fn(f.Name(), dat)
	}
}

func write(fp string, logger *log.Logger, v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Fatalf("%s: %v", fp, err)
	}

	// Make sure all directories exist to host this file
	dir := path.Dir(path.Join(*output, fp))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Fatalf("%s: %v", dir, err)
	}

	err = os.WriteFile(path.Join(*output, fp), out, 0644)
	if err != nil {
		logger.Fatalf("%s: %v", fp, err)
	}
}

func main() {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "meh: ", log.Llongfile)
	)

	input = flag.String("in", "", "path to the (uncompressed) medium archive")
	output = flag.String("out", "", "output directory")
	verbose = flag.Bool("verbose", false, "whether to print logs to stdout")
	flag.Parse()

	if *input == "" || *output == "" {
		flag.Usage()
		return
	}

	dirs, err := ioutil.ReadDir(*input)
	if err != nil {
		panic(err)
	}

	for _, d := range dirs {
		if d.IsDir() == false {
			logger.Printf("skipped %s: not a directory", d.Name())
			continue
		}

		switch d.Name() {
		case "blocks":
			users := []schema.User{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseBlocked(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				users = append(users, part...)
			})

			write("blocks.json", logger, schema.BlockedUsers{
				Meta:  "Blocked users",
				Users: users,
			})
		case "bookmarks":
			posts := []schema.Post{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseBookmarks(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				posts = append(posts, part...)
			})

			write("bookmarks.json", logger, schema.Bookmarks{
				Meta:  "Bookmarked posts",
				Posts: posts,
			})
		case "claps":
			claps := []schema.Clap{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseClaps(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				claps = append(claps, part...)
			})

			write("claps.json", logger, schema.Claps{
				Meta:  "Posts you've clapped for",
				Claps: claps,
			})
		case "interests":
			interests := schema.Interests{
				Meta: "Topics you're interested in",
			}

			walk(d.Name(), logger, func(name string, dat io.Reader) {
				switch name {
				case "publications.html":
					pubs, err := parser.ParseInterestsPublications(dat)
					if err != nil {
						logger.Fatalf("%s: %v", name, err)
						return
					}
					interests.Publications = pubs
				case "tags.html":
					tags, err := parser.ParseInterestsTags(dat)
					if err != nil {
						logger.Fatalf("%s: %v", name, err)
						return
					}
					interests.Tags = tags
				case "topics.html":
					topics, err := parser.ParseInterestsTopics(dat)
					if err != nil {
						logger.Fatalf("%s: %v", name, err)
						return
					}
					interests.Topics = topics
				case "writers.html":
					writers, err := parser.ParseInterestsWriters(dat)
					if err != nil {
						logger.Fatalf("%s: %v", name, err)
						return
					}
					interests.Writers = writers
				default:
					logger.Printf("Unknown interests file: %s", name)
				}
			})

			write("interests.json", logger, interests)
		case "ips":
			ips := []schema.IP{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseIps(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				ips = append(ips, part...)
			})

			write("ips.json", logger, schema.IPs{
				Meta: "Your IP history (note: Medium deletes IP history after 30 days)",
				IPs:  ips,
			})
		case "posts":
			posts := map[string]schema.Post{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				post, err := parser.ParsePost(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				posts[strings.TrimSuffix(name, ".html")] = *post
			})

			for name, post := range posts {
				write(path.Join("posts", name+".json"), logger, post)
			}
		case "lists":
			lists := []schema.List{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				list, err := parser.ParseList(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				lists = append(lists, *list)
			})

			write("lists.json", logger, schema.Lists{
				Meta:  "Lists you've created",
				Lists: lists,
			})
		case "pubs-following":
			pubs := []schema.Publication{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParsePublicationFollowing(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				pubs = append(pubs, part...)
			})

			write("following/publications.json", logger, schema.Publications{
				Meta:         "Publications you follow",
				Publications: pubs,
			})
		case "topics-following":
			topics := []schema.Topic{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseTopicsFollowing(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				topics = append(topics, part...)
			})

			write("following/topics.json", logger, schema.Topics{
				Meta:   "Topics you follow",
				Topics: topics,
			})
		case "users-following":
			users := []schema.User{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseUsersFollowing(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				users = append(users, part...)
			})

			write("following/users.json", logger, schema.Users{
				Meta:  "Users you follow",
				Users: users,
			})
		case "twitter":
			users := []schema.User{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseUsersSuggested(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				users = append(users, part...)
			})

			write("following/suggested.json", logger, schema.Users{
				Meta:  "Your Twitter friends who are also on Medium",
				Users: users,
			})
		case "sessions":
			sessions := []schema.Session{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseSessions(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				sessions = append(sessions, part...)
			})

			write("sessions.json", logger, schema.Sessions{
				Meta:     "Your active and inactive sessions across devices",
				Sessions: sessions,
			})
		case "highlights":
			highlights := []schema.Highlight{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseHighlights(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				highlights = append(highlights, part...)
			})

			write("highlights.json", logger, schema.Highlights{
				Meta:       "Your highlights",
				Highlights: highlights,
			})
		case "profile":
			profile := schema.Profile{
				Meta: "Your user profile",
			}
			profile.User = &schema.User{}

			walk(d.Name(), logger, func(name string, dat io.Reader) {
				switch {
				case name == "about.html":
					bio, _ := parser.ParseBio(dat)
					profile.User.Bio = bio
				case name == "profile.html":
					parser.ParseUserProfile(dat, &profile)
				case name == "publications.html":
					parser.ParsePublications(dat, &profile)
				case name == "memberships.html":
					parser.ParseMemberships(dat, &profile)
				case strings.HasPrefix(name, "charges-") && strings.HasSuffix(name, ".html"):
					parser.ParseMembershipCharges(dat, &profile)
				default:
					logger.Printf("skipped profile/%s: not supported", name)
				}
			})

			write("profile.json", logger, profile)
		default:
			logger.Printf("skipped %s: not supported", d.Name())
		}
	}

	if *verbose == true {
		fmt.Print(&buf)
	}
}
