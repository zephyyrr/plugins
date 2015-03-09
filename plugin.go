package plugins

type Plugin interface {
	Name() string
	Provides() []Event
	Subscribes() []Event

	Send(Event, Args) error
	Recieve() (Event, Args, error)
}

type Event string
type Args map[string]interface{}
