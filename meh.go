package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/valueof/meh/parser"
	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

var input *string
var output *string
var verbose *bool
var withImages *bool
var logger *log.Logger
var logbuf bytes.Buffer

func init() {
	input = flag.String("in", "", "path to the (uncompressed) medium archive")
	output = flag.String("out", "", "output directory")
	verbose = flag.Bool("verbose", false, "whether to print logs to stdout")
	withImages = flag.Bool("withImages", false, "whether to download images from medium cdn")
	logger = log.New(&logbuf, "meh: ", log.Lmsgprefix)
}

func walk(dir string, fn func(string, io.Reader)) error {
	dir = path.Join(*input, dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Printf("can't read %s, skipping", path.Base(dir))
		return err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") == false {
			logger.Printf("%s is not an html file, skipping", f.Name())
			continue
		}

		dat, err := os.Open(path.Join(dir, f.Name()))
		if err != nil {
			logger.Printf("can't read %s, skipping", f.Name())
			continue
		}
		defer dat.Close()

		fn(f.Name(), dat)
	}

	return nil
}

func write(fp string, v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Printf("can't marshal output for %s, skipping", fp)
		return
	}

	// Make sure all directories exist to host this file
	dir := path.Dir(path.Join(*output, fp))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Printf("can't create %s", dir)
		return
	}

	dest := path.Join(*output, fp)
	err = os.WriteFile(dest, out, 0644)
	if err != nil {
		logger.Printf("can't write to %s", dest)
	}
}

func download(img, dest string) {
	src := "https://cdn-images-1.medium.com/" + img
	out, err := os.Create(dest)
	if err != nil {
		logger.Printf("error creating file for %s: %v", img, err)
		return
	}
	defer out.Close()

	resp, err := http.Get(src)
	if err != nil {
		logger.Printf("error downloading %s: %v", src, err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Printf("error saving %s: %v", img, err)
		return
	}

	logger.Printf("downloaded %s", img)
}

func downloadImages(images []string) {
	dir := path.Join(*output, "images")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Printf("couldn't create %s, images will not be downloaded", dir)
		return
	}

	var wg sync.WaitGroup
	for _, img := range images {
		wg.Add(1)
		img := img
		go func() {
			defer wg.Done()
			download(img, path.Join(dir, img))
		}()
	}
	wg.Wait()
}

func main() {
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
			logger.Printf("%s is not a directory, skipping", d.Name())
			continue
		}

		switch d.Name() {
		case "blocks":
			users := []schema.User{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseBlocked(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				users = append(users, part...)
			})

			if err != nil {
				write("blocks.json", schema.BlockedUsers{
					Meta:  "Blocked users",
					Users: users,
				})
			}
		case "bookmarks":
			posts := []schema.Post{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseBookmarks(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				posts = append(posts, part...)
			})

			if err != nil {
				write("bookmarks.json", schema.Bookmarks{
					Meta:  "Bookmarked posts",
					Posts: posts,
				})
			}
		case "claps":
			claps := []schema.Clap{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseClaps(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				claps = append(claps, part...)
			})

			if err != nil {
				write("claps.json", schema.Claps{
					Meta:  "Posts you've clapped for",
					Claps: claps,
				})
			}
		case "interests":
			interests := schema.Interests{
				Meta: "Topics you're interested in",
			}

			err = walk(d.Name(), func(name string, dat io.Reader) {
				switch name {
				case "publications.html":
					pubs, err := parser.ParseInterestsPublications(dat)
					if err != nil {
						logger.Printf("error parsing %s, skipping", name)
						return
					}
					logger.Printf("parsed %s", name)
					interests.Publications = pubs
				case "tags.html":
					tags, err := parser.ParseInterestsTags(dat)
					if err != nil {
						logger.Printf("error parsing %s, skipping", name)
						return
					}
					logger.Printf("parsed %s", name)
					interests.Tags = tags
				case "topics.html":
					topics, err := parser.ParseInterestsTopics(dat)
					if err != nil {
						logger.Printf("error parsing %s, skipping", name)
						return
					}
					logger.Printf("parsed %s", name)
					interests.Topics = topics
				case "writers.html":
					writers, err := parser.ParseInterestsWriters(dat)
					if err != nil {
						logger.Printf("error parsing %s, skipping", name)
						return
					}
					logger.Printf("parsed %s", name)
					interests.Writers = writers
				default:
					logger.Printf("Unknown interests file %s, skipping", name)
				}
			})

			if err != nil {
				write("interests.json", interests)
			}
		case "ips":
			ips := []schema.IP{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseIps(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				ips = append(ips, part...)
			})

			if err != nil {
				write("ips.json", schema.IPs{
					Meta: "Your IP history (note: Medium deletes IP history after 30 days)",
					IPs:  ips,
				})
			}
		case "posts":
			posts := map[string]schema.Post{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				post, err := parser.ParsePost(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				posts[strings.TrimSuffix(name, ".html")] = *post
			})

			if err != nil {
				for name, post := range posts {
					write(path.Join("posts", name+".json"), post)
				}
			}
		case "lists":
			lists := []schema.List{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				list, err := parser.ParseList(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				lists = append(lists, *list)
			})

			if err != nil {
				write("lists.json", schema.Lists{
					Meta:  "Lists you've created",
					Lists: lists,
				})
			}
		case "pubs-following":
			pubs := []schema.Publication{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParsePublicationFollowing(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				pubs = append(pubs, part...)
			})

			if err != nil {
				write("following/publications.json", schema.Publications{
					Meta:         "Publications you follow",
					Publications: pubs,
				})
			}
		case "topics-following":
			topics := []schema.Topic{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseTopicsFollowing(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				topics = append(topics, part...)
			})

			if err != nil {
				write("following/topics.json", schema.Topics{
					Meta:   "Topics you follow",
					Topics: topics,
				})
			}
		case "users-following":
			users := []schema.User{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseUsersFollowing(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				users = append(users, part...)
			})

			if err != nil {
				write("following/users.json", schema.Users{
					Meta:  "Users you follow",
					Users: users,
				})
			}
		case "twitter":
			users := []schema.User{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseUsersSuggested(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				users = append(users, part...)
			})

			if err != nil {
				write("following/suggested.json", schema.Users{
					Meta:  "Your Twitter friends who are also on Medium",
					Users: users,
				})
			}
		case "sessions":
			sessions := []schema.Session{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseSessions(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				sessions = append(sessions, part...)
			})

			if err != nil {
				write("sessions.json", schema.Sessions{
					Meta:     "Your active and inactive sessions across devices",
					Sessions: sessions,
				})
			}
		case "highlights":
			highlights := []schema.Highlight{}
			err = walk(d.Name(), func(name string, dat io.Reader) {
				part, err := parser.ParseHighlights(dat)
				if err != nil {
					logger.Printf("error parsing %s, skipping", name)
					return
				}
				logger.Printf("parsed %s", name)
				highlights = append(highlights, part...)
			})

			if err != nil {
				write("highlights.json", schema.Highlights{
					Meta:       "Your highlights",
					Highlights: highlights,
				})
			}
		case "profile":
			profile := schema.Profile{
				Meta: "Your user profile",
			}
			profile.User = &schema.User{}

			err = walk(d.Name(), func(name string, dat io.Reader) {
				switch {
				case name == "about.html":
					bio, err := parser.ParseBio(dat)
					if err != nil {
						logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					logger.Printf("parsed %s", name)
					profile.User.Bio = bio
				case name == "profile.html":
					err = parser.ParseUserProfile(dat, &profile)
					if err != nil {
						logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					logger.Printf("parsed %s", name)
				case name == "publications.html":
					err = parser.ParsePublications(dat, &profile)
					if err != nil {
						logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					logger.Printf("parsed %s", name)
				case name == "memberships.html":
					err = parser.ParseMemberships(dat, &profile)
					if err != nil {
						logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					logger.Printf("parsed %s", name)
				case strings.HasPrefix(name, "charges-") && strings.HasSuffix(name, ".html"):
					err = parser.ParseMembershipCharges(dat, &profile)
					if err != nil {
						logger.Printf("error parsing %s, profile.json will be incomplete", name)
						return
					}
					logger.Printf("parsed %s", name)
				default:
					logger.Printf("skipped profile/%s: not supported", name)
				}
			})

			if err != nil {
				write("profile.json", profile)
			}
		default:
			logger.Printf("%s isn't supported, skipping", d.Name())
		}
	}

	if *withImages {
		downloadImages(util.GetQueuedImages())
	} else {
		logger.Printf("not downloading images, use -withImages if you want to download images")
	}

	if *verbose {
		fmt.Print(&logbuf)
	}
}
