package plugins

import (
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
)

var logger = log.New(os.Stderr, "[plugins::Manager] ", log.LstdFlags)

type Manager struct {
	plugins       plugins
	closers       map[string]chan<- struct{}
	subscriptions map[Event]plugins
	selects       []reflect.SelectCase
	handler       Handler
}

type plugins map[string]Plugin

// Returns a new plugin Manager ready for use.
// A default Muxer is installed as the handler
func NewManager() *Manager {
	return &Manager{
		plugins:       make(plugins),
		closers:       make(map[string]chan<- struct{}),
		subscriptions: make(map[Event]plugins),
		selects:       make([]reflect.SelectCase, 0, 16),
		handler:       make(mapMuxr),
	}
}

type packet struct {
	Event Event
	Args  Args
}

func (ph *Manager) Handle(pl Plugin) {
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

func (ph *Manager) HandleAll(pls ...Plugin) {
	for _, pl := range pls {
		ph.Handle(pl)
	}
}

func (ph Manager) Handler() Handler {
	return ph.handler
}

func (ph Manager) SetHandler(h Handler) {
	ph.handler = h
}

//Returns the installed handler if it is a Muxer
//Otherwise, returns nil
func (ph Manager) Muxer() Muxer {
	if muxr, ok := ph.handler.(Muxer); ok {
		return muxr
	}
	return nil
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
						fallthrough
					default:
						logger.Printf("Encountered a error recieving event from %s: \n%s", pl.Name(), err.Error())
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

func (ph *Manager) Unload(pl Plugin) {
	close(ph.closers[pl.Name()])
	delete(ph.closers, pl.Name())

	for _, sub := range pl.Subscribes() {
		delete(ph.subscriptions[sub], pl.Name())
	}

	delete(ph.plugins, pl.Name())
}

func (ph *Manager) ListenAndServe() error {
	for {
		if len(ph.selects) <= 0 {
			runtime.Gosched()
			continue
		}
		if chosen, recv, ok := reflect.Select(ph.selects); ok {
			pck := recv.Interface().(packet)
			ph.dispatch("plugin", pck.Event, pck.Args)
		} else if chosen >= 0 && chosen < len(ph.selects) {
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

func (ph Manager) Dispatch(e Event, args Args) error {
	return ph.dispatch("local", e, args)
}

func (ph Manager) dispatch(identifier string, e Event, args Args) (err error) {
	ph.Handler().HandleEvent(e, args)
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
