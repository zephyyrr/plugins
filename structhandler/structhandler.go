package structhandler

import (
	"github.com/mitchellh/mapstructure"
	"github.com/zephyyrr/plugins"
	"reflect"
)

type typedHandler struct {
	f reflect.Value
	t reflect.Type
}

type untypedHandler struct {
	f reflect.Value
}

//Creates a new Handler that unpackes the Args into the argument of f, that has to be a function taking one parameter
//of kind struct or no parameters.
//When the resulting handler is called, it will call f using the unpacked arguments.
func New(f interface{}) plugins.Handler {
	ft := reflect.TypeOf(f)
	if ft.NumIn() == 0 {
		return untypedHandler{reflect.ValueOf(f)}
	}
	return typedHandler{reflect.ValueOf(f), ft.In(0)}
}

func (h typedHandler) HandleEvent(e plugins.Event, args plugins.Args) {
	s := reflect.New(h.t)
	err := mapstructure.Decode(args, s.Interface())
	if err != nil {
		panic(err)
	}
	h.f.Call([]reflect.Value{reflect.Indirect(s)})
}

func (h untypedHandler) HandleEvent(e plugins.Event, args plugins.Args) {
	h.f.Call(nil)
}
