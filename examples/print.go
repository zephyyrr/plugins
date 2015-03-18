package main

import (
	"fmt"
	"github.com/zephyyrr/plugins"
)

type PrintPlugin struct{}

func (PrintPlugin) Name() string {
	return "Print"
}

func (PrintPlugin) Provides() []plugins.Event {
	return nil
}

func (PrintPlugin) Subscribes() []plugins.Event {
	return []plugins.Event{
		"message",
	}
}

func (PrintPlugin) Send(e plugins.Event, args plugins.Args) error {
	if e != "message" {
		panic("Wrong event")
	}

	fmt.Println(args["message"])
	return nil
}

func (PrintPlugin) Recieve() (plugins.Event, plugins.Args, error) {
	return "", nil, plugins.Unblocking
}
