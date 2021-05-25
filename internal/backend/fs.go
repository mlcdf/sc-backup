package backend

import (
	"os"
	"path"

	"go.mlcdf.fr/sc-backup/internal/domain"
)

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ domain.Backend = (*fs)(nil)

type fs struct {
	location  string
	formatter domain.Formatter
}

func NewFS(location string, format domain.Formatter) *fs {
	return &fs{location, format}
}

func (f *fs) Create() error {
	os.MkdirAll(f.location, os.ModePerm)
	return nil
}

func (f *fs) Location() string {
	return f.location
}

func (f *fs) Save(data domain.Serializable) error {
	p := path.Join(f.location, data.Slug()+f.formatter.Ext())

	fd, err := os.Create(p)
	if err != nil {
		return err
	}

	return f.formatter.Format(data, fd)
}
