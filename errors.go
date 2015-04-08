package plugins

type Error int

func (e Error) Error() string {
	if s, ok := errors[e]; ok {
		return s
	}

	if e < 100 {
		return "Unknown transmission error"
	} else if e < 200 {
		return "Unknown server error"
	} else if e < 300 {
		return "Unknown client error"
	} else {
		return "Unknown error"
	}
}

var errors = map[Error]string{
	Success:           "No error. Everything is fine.",
	NoSupportedFormat: "No supported formats in common",
	Unblocking:        "This is a unblocking recieve",
	NotImplemented:    "Operation is not implemented",
	NotDirectory:      "Path supplied is not a directory",
	DuplicatePlugin:   "Plugin is a duplicate. Already handled.",
}

//Transmission errors
const (
	Success Error = iota
	NoSupportedFormat
)

//Server errors
const (
	Unblocking Error = 100 + iota
	NotImplemented
	NotDirectory
	DuplicatePlugin
)
