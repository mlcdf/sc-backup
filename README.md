# sc-backup

/!\ DOES NOT WORK SINCE THE LATEST SENSCRITIQUE UPDATE WITH USES CLIENT-SIDE RENDERING :/

/!\ I since changed the URL to old.senscritique.com, so this will actually work (as long as old.senscritique.com refers to the old server-side-rendered website)

[![test](https://github.com/mlcdf/sc-backup/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/mlcdf/sc-backup/actions/workflows/test.yml)
[![coverage](https://mlcdf.github.io/sc-backup/badge.svg)](https://mlcdf.github.io/sc-backup/report.html)

A fast and easy way to backup a [SensCritique](https://www.senscritique.com) user or list.

## Install

- From [GitHub releases](https://github.com/mlcdf/sc-backup/releases): download the binary corresponding to your OS and architecture.
- From source (make sure `$GOPATH/bin` is in your `$PATH`):
```sh
go get go.mlcdf.fr/sc-backup
```

## Usage

```
Usage:
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
```

Check out the [examples](examples) to see what the output looks like.

## Development

Run the app
```sh
go run main.go
```

Run the tests
```sh
go test ./...
```
