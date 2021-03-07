// +build !release

package backend

// Memory is used for testing purpose
type Memory struct {
	Stuff interface{}
	Slug  string
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Create() error {
	return nil
}

func (m *Memory) Location() string {
	return "memory"
}

func (m *Memory) SaveCollection(stuff interface{}, slug string) error {
	m.Stuff = stuff
	m.Slug = slug
	return nil
}

func (m *Memory) SaveList(stuff interface{}, slug string) error {
	m.Stuff = stuff
	m.Slug = slug
	return nil
}
