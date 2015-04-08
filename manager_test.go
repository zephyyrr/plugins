package plugins

import (
	"reflect"
	"testing"
)

func TestManager_Handle(t *testing.T) {
	man := NewManager()
	man.Handle(TestPlugin{})

	if _, ok := man.plugins["TestPlugin"]; !ok {
		t.Fail()
	}

	if _, ok := man.subscriptions["test.answer"]; !ok {
		t.Fail()
	}
}

func TestManager_Handler(t *testing.T) {
	man := NewManager()

	if man.Handler() == nil {
		t.Fail()
	}

	if _, ok := man.Handler().(Muxer); !ok {
		t.Fail()
	}

	if !reflect.DeepEqual(reflect.ValueOf(man.handler), reflect.ValueOf(man.Handler())) {
		t.Fail()
	}
}

func TestManager_SetHandler(t *testing.T) {
	man := NewManager()

	man.SetHandler(nil)

	if man.Handler() != nil {
		//t.Fail()
	}

	var h Handler = make(mapMuxr)

	man.SetHandler(h)

	if !reflect.DeepEqual(reflect.ValueOf(h), reflect.ValueOf(man.Handler())) {
		t.Fail()
	}
}

func TestManager_Unhandle(t *testing.T) {

}
