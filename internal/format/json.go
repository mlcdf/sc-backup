package format

import (
	"encoding/json"
	"io"

	"go.mlcdf.fr/sc-backup/internal/domain"
)

var _ domain.Formatter = (*JSON)(nil)

type JSON struct {
	pretty bool
}

func NewJSON(pretty bool) *JSON {
	return &JSON{pretty}
}

func (f *JSON) Ext() string {
	return ".json"
}

func (f *JSON) Format(data domain.Serializable, writer io.Writer) error {
	var formatted []byte
	var err error

	if f.pretty {
		formatted, err = json.MarshalIndent(data.JSON(), "", "    ")
	} else {
		formatted, err = json.Marshal(data.JSON())
	}
	if err != nil {
		return err
	}

	_, err = writer.Write(formatted)
	return err
}
