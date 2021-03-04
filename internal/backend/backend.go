package backend

type Backend interface {
	// Location returns this backend's location (the directory name).
	Location() string

	// Create the backend resources
	Create() error

	SaveCollection(stuff interface{}, slug string) error
	SaveList(stuff interface{}, slug string) error
}
