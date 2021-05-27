// +build !release

package mock

import "go.mlcdf.fr/sc-backup/internal/domain"

var _ domain.Backend = (*Backend)(nil)

// Backend is used for testing purpose
type Backend struct {
	Data map[string]interface{}
}

func NewBackend() *Backend {
	return &Backend{
		Data: map[string]interface{}{},
	}
}

func (m *Backend) Create() error {
	return nil
}

func (m *Backend) Location() string {
	return "Mock"
}

func (m *Backend) Save(data domain.Serializable) error {
	m.Data[data.Slug()] = data
	return nil
}
