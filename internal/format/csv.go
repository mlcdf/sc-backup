package format

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"go.mlcdf.fr/sc-backup/internal/domain"
)

var _ domain.Formatter = (*CSV)(nil)

type CSV struct{}

func (f *CSV) Ext() string {
	return ".csv"
}

func (f *CSV) Format(data domain.Serializable, writer io.Writer) error {
	mapMapString := make([][]string, 0, len(data.CSV()))

	w := csv.NewWriter(writer)
	for _, entry := range data.CSV() {
		mapString := []string{
			entry.ID,
			entry.Title,
			entry.OriginalTitle,
			strconv.Itoa(entry.Year),
			strings.Join(entry.Authors, ";"),
			strconv.Itoa(entry.Rating),
		}
		mapMapString = append(mapMapString, mapString)
	}

	err := w.WriteAll(mapMapString)
	return err
}
