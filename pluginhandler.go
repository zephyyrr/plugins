package plugins

import (
	"reflect"
	"runtime"
	"time"
)

type Handler struct {
	plugins       plugins
	closers       map[string]chan<- struct{}
	subscriptions map[Event]plugins
	selects       []reflect.SelectCase
}

type plugins map[string]Plugin

func NewHandler() *Handler {
	return &Handler{
		plugins:       make(plugins),
		closers:       make(map[string]chan<- struct{}),
		subscriptions: make(map[Event]plugins),
		selects:       make([]reflect.SelectCase, 0, 16),
	}
}

type packet struct {
	Event Event
	Args  Args
}

func (ph *Handler) Handle(pl Plugin) {
	ph.plugins[pl.Name()] = pl

	for _, sub := range pl.Subscribes() {
		if ph.subscriptions[sub] == nil {
			ph.subscriptions[sub] = make(plugins)
		}
		ph.subscriptions[sub][pl.Name()] = pl
	}

	f, ch, closer := generatefunc(pl)
	ph.closers[pl.Name()] = closer

	go f()

	rch := reflect.ValueOf(ch)
	ph.selects = append(ph.selects,
		reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: rch,
		})

}

func (ph *Handler) HandleAll(pls ...Plugin) {
	for _, pl := range pls {
		ph.Handle(pl)
	}
}

func generatefunc(pl Plugin) (func(), <-chan packet, chan<- struct{}) {
	ch, closer := make(chan packet), make(chan struct{})
	return func() {
		defer close(ch)

		if len(pl.Provides()) == 0 {
			return //Nothing is provided, so nothing to listen for.
		}

		for {
			select {
			case <-closer:
				return
			default:
				event, args, err := pl.Recieve() //Probably blocking. Hopefully not for to long.
				if err != nil {
					switch err {
					case Unblocking:
						time.Sleep(500 * time.Millisecond)
					case NotImplemented:
					default:
						return
					}
					continue
				}
				ch <- packet{event, args}
			}
			runtime.Gosched() // Allow something else to run if tight loop.
		}
	}, ch, closer
}

func (ph *Handler) Unload(pl Plugin) {
	close(ph.closers[pl.Name()])
	delete(ph.closers, pl.Name())

	for _, sub := range pl.Subscribes() {
		delete(ph.subscriptions[sub], pl.Name())
	}

	delete(ph.plugins, pl.Name())
}

func (ph *Handler) ListenAndServe() error {

	for len(ph.plugins) > 0 {
		if chosen, recv, ok := reflect.Select(ph.selects); ok {
			pck := recv.Interface().(packet)
			ph.Dispatch(pck.Event, pck.Args)
		} else {
			//Not ok
			//Removing select case from list
			if chosen == len(ph.selects)-1 {
				ph.selects = ph.selects[:len(ph.selects)-1]
			} else {
				ph.selects = append(ph.selects[:chosen], ph.selects[chosen+1:]...)
			}
		}
	}

	return nil
}

func (ph Handler) Dispatch(e Event, args Args) error {
	return ph.dispatch("local", e, args)
}

func (ph Handler) dispatch(identifier string, e Event, args Args) (err error) {
	for _, plugin := range ph.subscriptions[e] {
		go func(plugin Plugin) {
			tmp := plugin.Send(e, args)
			if err == nil && tmp != NotImplemented {
				err = tmp
			}
		}(plugin)
	}
	return
}
