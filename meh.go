package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/valueof/meh/parser"
)

func walk(dir string, logger *log.Logger, fn func(string, io.Reader)) {
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
		if err != nil {
			logger.Fatalf("%s: %v", f.Name(), err)
			continue
		}

		fn(f.Name(), dat)
	}
}

func write(name string, logger *log.Logger, v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Fatalf("%s: %v", name, err)
	}

	err = os.WriteFile("./out/"+name+".json", out, 0644)
	if err != nil {
		logger.Fatalf("%s: %v", name, err)
	}
}

func main() {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "meh: ", log.Llongfile)
		root   = "./data"
	)

	dirs, err := ioutil.ReadDir(root)
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
			users := []parser.User{}
			walk(path.Join(root, d.Name()), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseBlocked(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				users = append(users, part...)
			})

			write("blocks", logger, parser.BlockedUsers{
				Users: users,
			})
		case "bookmarks":
			posts := []parser.Post{}
			walk(path.Join(root, d.Name()), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseBookmarks(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				posts = append(posts, part...)
			})

			write("bookmarks", logger, parser.Bookmarks{
				Posts: posts,
			})
		case "claps":
			claps := []parser.Clap{}
			walk(path.Join(root, d.Name()), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseClaps(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				claps = append(claps, part...)
			})

			write("claps", logger, parser.Claps{
				Claps: claps,
			})
		case "interests":
			interests := parser.Interests{}
			walk(path.Join(root, d.Name()), logger, func(name string, dat io.Reader) {
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

			write("interests", logger, interests)
		}
	}

	fmt.Print(&buf)
}
