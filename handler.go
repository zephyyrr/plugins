package plugins

type Handler interface {
	HandleEvent(Event, Args)
}

type HandlerFunc func(Event, Args)

func (hf HandlerFunc) HandleEvent(event Event, args Args) {
	hf(event, args)
}

type Muxer interface {
	Handler
	AddHandler(Event, Handler)
	RemoveHandler(Event, Handler)
}

type mapMuxr map[Event]Handler

func (m mapMuxr) HandleEvent(event Event, args Args) {
	if handler, ok := m[event]; ok {
		handler.HandleEvent(event, args)
	}
}

func (m mapMuxr) AddHandler(event Event, handlr Handler) {
	m[event] = handlr
}

func (m mapMuxr) RemoveHandler(event Event, handlr Handler) {
	if h, ok := m[event]; ok && h == handlr {
		delete(m, event)
	}
}
