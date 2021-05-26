package domain

type Serializable interface {
	Slug() string
	CSV() []*Entry
	JSON() interface{}
}

type Backend interface {
	// Location returns this backend's location (the directory name).
	Location() string

	// Create the backend resources
	Create() error

	Save(Serializable) error
}
