package plugins

import (
	"errors"
)

var (
	NoSupportedFormat = errors.New("No supported formats in common")
	Unblocking        = errors.New("This is a unblocking recieve")
	NotImplemented    = errors.New("Operation is not implemented")
	NotDirectory      = errors.New("Path supplied is not a directory")
)
