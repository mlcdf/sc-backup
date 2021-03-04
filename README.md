# sc-backup

A fast and easy way to backup a [SensCritique](https://senscritique.com) user or list.

## Installation

- From [GitHub releases](https://github.com/mlcdf/sc-backup/releases): download and place the binary in your `$PATH`
- From source (first, make sure you have `$GOPATH/bin` in your `$PATH`):
```sh
go get https://github.com/mlcdf/sc-backup
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

## Development

Run the app
```sh
go run main.go
```

Lancer les tests
```sh
go test ./...
```
