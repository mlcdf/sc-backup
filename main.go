package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/mlcdf/sc-backup/internal/backend"
	"github.com/mlcdf/sc-backup/internal/backup"
	"github.com/mlcdf/sc-backup/internal/logx"
)

const usage = `Usage:
    sc-backup --collection [USERNAME]
    sc-backup --list [URL]

Options:
    -c, --collection USERNAME   Backup a user's collection
    -l, --list URL              Backup a list
    -o, --output PATH           Directory at which to backup the data. Defaults to ./output
    -f, --format json|csv       Export format. Defaults to json
    -p, --pretty                Prettify the JSON exports
    -v, --verbose               Print verbose output
    -V, --version               Print version

Examples:
    sc-backup --collection mlcdf
    sc-backup --list https://www.senscritique.com/liste/Vu_au_cinema/363578
`

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module. See
// golang.org/issue/29814 and golang.org/issue/29228.
var Version string

func main() {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	var (
		isVerboseFlag  bool
		listFlag       string
		collectionFlag string
		outputFlag     string = "output"
		formatFlag     string = "json"
		prettyFlag     bool
		versionFlag    bool
	)

	flag.BoolVar(&versionFlag, "version", versionFlag, "print the version")
	flag.BoolVar(&versionFlag, "V", versionFlag, "print the version")

	flag.BoolVar(&isVerboseFlag, "verbose", isVerboseFlag, "enable verbose output")
	flag.BoolVar(&isVerboseFlag, "v", isVerboseFlag, "enable verbose output")

	flag.StringVar(&listFlag, "list", listFlag, "Download list")
	flag.StringVar(&listFlag, "l", listFlag, "Download list")

	flag.StringVar(&collectionFlag, "collection", collectionFlag, "Download user collection")
	flag.StringVar(&collectionFlag, "-c", collectionFlag, "Download user collection")

	flag.StringVar(&outputFlag, "output", outputFlag, "Output directory")
	flag.StringVar(&outputFlag, "o", outputFlag, "Output directory")

	flag.StringVar(&formatFlag, "format", formatFlag, "Output format. Either json or csv. Default to json.")
	flag.StringVar(&formatFlag, "f", formatFlag, "Output format. Either json or csv. Default to json.")

	flag.BoolVar(&prettyFlag, "pretty", prettyFlag, "Pretty output")
	flag.BoolVar(&prettyFlag, "p", prettyFlag, "Pretty output")

	flag.Parse()

	if versionFlag {
		if Version != "" {
			fmt.Println(Version)
			return
		}
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(buildInfo.Main.Version)
			return
		}
		fmt.Println("(unknown)")
		return
	}

	start := time.Now()

	if collectionFlag != "" && listFlag != "" {
		log.Fatalln("error: you can't set --list and --collection at the same time")
	}

	if collectionFlag == "" && listFlag == "" {
		log.Fatalln("error: at least one of --list or --collection is required")
	}

	if isVerboseFlag == true {
		logx.EnableVerboseOutput()
	}

	var back backend.Backend
	var err error

	if collectionFlag != "" {
		back, err = backend.NewFS(filepath.Join(outputFlag, collectionFlag), prettyFlag, formatFlag)
		if err != nil {
			log.Fatalf("error: %s", err)
		}

		err = backup.Collection(collectionFlag, back)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	if listFlag != "" {
		back, err = backend.NewFS(outputFlag, prettyFlag, formatFlag)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		err = backup.List(listFlag, back)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	to, err := filepath.Abs(back.Location())
	if err != nil {
		to = back.Location()
	}
	logx.Info("Saved to %s in %s", to, time.Since(start).Round(time.Millisecond).String())
}
