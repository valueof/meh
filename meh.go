package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/valueof/meh/formatters"
	"github.com/valueof/meh/parser"
	http "github.com/valueof/meh/server"
	"github.com/valueof/meh/util"
)

var VERSION string = "0.3"

var dir *string
var zip *string
var output *string
var verbose *bool
var withImages *bool
var version *bool
var server *string
var logger *log.Logger
var logbuf bytes.Buffer

func init() {
	dir = flag.String("dir", "", "path to the uncompressed medium archive")
	zip = flag.String("zip", "", "path to the compressed medium archive")
	output = flag.String("out", "", "output directory")
	server = flag.String("server", "", "run web version of meh on provided address")
	verbose = flag.Bool("verbose", false, "whether to print logs to stdout")
	version = flag.Bool("version", false, "print version and exit")
	withImages = flag.Bool("withImages", false, "whether to download images from medium cdn")
	logger = log.New(&logbuf, "meh: ", log.Lmsgprefix)
}

func run() error {
	if *verbose {
		logger.SetOutput(os.Stdout)
	}

	if *version {
		fmt.Printf("meh %s\n", VERSION)
		return nil
	}

	if *server != "" {
		http.RunHTTPServer(*server)
		return nil
	}

	if (*dir == "" && *zip == "") || *output == "" {
		flag.Usage()
		return nil
	}

	input := ""

	switch {
	case *zip != "":
		tmp := filepath.Join(*output, ".archive")
		err := util.UnzipArchive(*zip, tmp)
		if err != nil {
			logger.Printf("UnzipArchive(%s, %s): %v", *zip, tmp, err)
			return err
		}

		defer func() {
			logger.Printf("clean up: removing %s", tmp)
			os.RemoveAll(tmp)
		}()

		logger.Printf("extracted archive into %s", tmp)
		input, err = util.FindArchiveRoot(tmp)
		if err != nil {
			logger.Printf("FindArchiveRoot(%s): %v", tmp, err)
			return err
		}
	case *dir != "":
		input = *dir
	}

	input, err := filepath.Abs(input)
	if err != nil {
		logger.Printf("filepath.Abs(): %v", err)
		return err
	}

	logger.Printf("using directory %s as input", input)

	w := formatters.NewJSONFormatter(*output, *logger)
	p := parser.NewParser(input, *logger, w)
	err = p.Parse()
	if err != nil {
		logger.Printf("parser.Parse(): %v", err)
		return err
	}

	if *withImages {
		p.FetchImages(*output)
	} else {
		logger.Printf("not downloading images, use -withImages if you want to download images")
	}

	return nil
}

func main() {
	flag.Parse()
	err := run()

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
