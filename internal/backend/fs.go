package backend

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mlcdf/sc-backup/internal/sc"
)

type fs struct {
	location string
	pretty   bool
	format   string
}

func NewFS(location string, pretty bool, format string) (*fs, error) {
	f := &fs{location: location, pretty: pretty}
	if format == "json" || format == "csv" {
		f.format = format
	} else {
		return nil, fmt.Errorf("invalid format %s: it should be either 'json' or 'csv'", format)
	}

	return f, nil
}

func (f *fs) Create() error {
	os.MkdirAll(filepath.Join(f.location), os.ModePerm)
	return nil
}

func (f *fs) Location() string {
	return f.location
}

func (f *fs) SaveCollection(stuff interface{}, slug string) error {
	if f.format == "json" {
		return writeJSON(stuff, filepath.Join(f.location, slug+".json"), f.pretty)
	} else {
		return writeCSV(stuff, filepath.Join(f.location, slug+".csv"), f.pretty)
	}
}

func (f *fs) SaveList(stuff interface{}, slug string) error {
	if f.format == "json" {
		return writeJSON(stuff, filepath.Join(f.location, slug+".json"), f.pretty)
	} else {
		return writeCSV(stuff, filepath.Join(f.location, slug+".csv"), f.pretty)
	}
}

func writeJSON(stuff interface{}, filename string, pretty bool) error {
	var jsonString []byte
	var err error

	if pretty {
		jsonString, err = json.MarshalIndent(stuff, "", "    ")
	} else {
		jsonString, err = json.Marshal(stuff)
	}
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, jsonString, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func writeCSV(entries interface{}, filename string, pretty bool) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	mapMapString := make([][]string, 0, len(entries.([]*sc.Entry)))

	writer := csv.NewWriter(f)
	for _, entry := range entries.([]*sc.Entry) {
		mapString := []string{
			entry.ID,
			entry.Title,
			entry.FrenchTitle,
			strconv.Itoa(entry.Year),
			strings.Join(entry.Authors, ";"),
			strconv.Itoa(entry.Rating),
		}
		mapMapString = append(mapMapString, mapString)
	}

	err = writer.WriteAll(mapMapString)
	if err != nil {
		return err
	}

	return nil
}
