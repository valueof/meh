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

func marshal(name string, logger *log.Logger, v any) {
	out, err := json.Marshal(v)
	if err != nil {
		logger.Fatalf("%s: %v", name, err)
	}
	fmt.Printf("%s: %s", name, string(out))
	fmt.Println()
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

			marshal("blocks", logger, parser.BlockedUsers{
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

			marshal("bookmarks", logger, parser.Bookmarks{
				Posts: posts,
			})
		}
	}

	// fmt.Print(&buf)
}
