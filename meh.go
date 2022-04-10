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

	dir := path.Join(*output, name+".json")
	err = os.WriteFile(dir, out, 0644)
	if err != nil {
		logger.Fatalf("%s: %v", name, err)
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
			users := []parser.User{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
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
			walk(d.Name(), logger, func(name string, dat io.Reader) {
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
			walk(d.Name(), logger, func(name string, dat io.Reader) {
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

			write("interests", logger, interests)
		case "ips":
			ips := []parser.IP{}
			walk(d.Name(), logger, func(name string, dat io.Reader) {
				part, err := parser.ParseIps(dat)
				if err != nil {
					logger.Fatalf("%s: %v", name, err)
					return
				}
				ips = append(ips, part...)
			})

			write("ips", logger, parser.IPs{
				IPs: ips,
			})
		default:
			logger.Printf("skipped %s: not supported", d.Name())
		}
	}

	if *verbose == true {
		fmt.Print(&buf)
	}
}
