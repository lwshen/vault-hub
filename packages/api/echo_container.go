package api

// Container holds all dependencies for the application handlers
type Container struct {
	// Add any dependencies here if needed
}

// NewContainer returns a new container with initialized dependencies
func NewContainer() (*Container, error) {
	c := &Container{}
	return c, nil
}
