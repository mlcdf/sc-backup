// +build !release

package backend

import "go.mlcdf.fr/sc-backup/internal/domain"

var _ domain.Backend = (*Memory)(nil)

// Memory is used for testing purpose
type Memory struct {
	Data map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{
		Data: map[string]interface{}{},
	}
}

func (m *Memory) Create() error {
	return nil
}

func (m *Memory) Location() string {
	return "memory"
}

func (m *Memory) Save(data domain.Serializable) error {
	m.Data[data.Slug()] = data
	return nil
}
