package domain

import "io"

type Formatter interface {
	// Format the entries and save them to the backend
	Format(entries Serializable, writer io.Writer) error
	// Ext returns the file extension
	Ext() string
}
