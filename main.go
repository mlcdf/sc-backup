package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"go.mlcdf.fr/sc-backup/internal/backend"
	"go.mlcdf.fr/sc-backup/internal/backup"
	"go.mlcdf.fr/sc-backup/internal/domain"
	"go.mlcdf.fr/sc-backup/internal/format"
	"go.mlcdf.fr/sc-backup/internal/logging"
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
	flag.StringVar(&collectionFlag, "c", collectionFlag, "Download user collection")

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

	if formatFlag == "csv" && prettyFlag {
		logging.Info("warning: -p/--pretty is useless with -f/--format csv. CSV won't be prettified.")
	}

	if isVerboseFlag {
		logging.EnableVerboseOutput()
	}

	var back domain.Backend
	var err error

	var formatter domain.Formatter

	switch formatFlag {
	case "json":
		formatter = format.NewJSON(prettyFlag)
	case "csv":
		formatter = &format.CSV{}
	default:
		log.Fatalf("invalid format %s: it should be json|csv|html", formatFlag)
	}

	if collectionFlag != "" {
		back = backend.NewFS(filepath.Join(outputFlag, collectionFlag), formatter)
		err = backup.Collection(collectionFlag, back)
	}

	if listFlag != "" {
		back = backend.NewFS(outputFlag, formatter)
		err = backup.List(listFlag, back)
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	to, err := filepath.Abs(back.Location())
	if err != nil {
		to = back.Location()
	}
	logging.Info("Saved to %s in %s", to, time.Since(start).Round(time.Millisecond).String())
}
