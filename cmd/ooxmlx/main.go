package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/example/ooxmlx/internal/buildinfo"
	"github.com/example/ooxmlx/internal/extract"
	"github.com/example/ooxmlx/internal/transform"
)

func main() {
	var dest string
	var overwrite bool
	var indent string
	var encoding string
	var fixNBSP bool
	var quiet bool
	var showVersion bool

	flag.StringVar(&dest, "o", "", "Destination directory (required)")
	flag.BoolVar(&overwrite, "overwrite", false, "Allow extraction into non-empty dest")
	flag.StringVar(&indent, "indent", "  ", "Indentation for XML output")
	flag.StringVar(&encoding, "encoding", "utf-8", "XML encoding declaration")
	flag.BoolVar(&fixNBSP, "fix-nbsp", false, "Normalize non-breaking spaces in XML")
	flag.BoolVar(&quiet, "quiet", false, "Suppress logs")
	flag.BoolVar(&showVersion, "version", false, "Print version info")
	flag.Parse()

	if showVersion {
		fmt.Println(buildinfo.Summary())
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: ooxmlx [flags] <archive>")
		flag.PrintDefaults()
		os.Exit(2)
	}
	if dest == "" {
		fmt.Fprintln(os.Stderr, "-o destination is required")
		os.Exit(2)
	}

	archivePath := args[0]

	var transformer transform.Transformer = transform.Nop{}
	if fixNBSP {
		transformer = transform.Composite{Transformers: []transform.Transformer{transform.ReplaceNBSP{}}}
	}

	var logger *log.Logger
	if !quiet {
		logger = log.New(os.Stdout, "ooxmlx ", log.LstdFlags)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err := extract.Extract(ctx, archivePath, extract.Options{
		Destination: dest,
		Overwrite:   overwrite,
		Indent:      indent,
		Encoding:    encoding,
		Transformer: transformer,
		Logger:      logger,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
