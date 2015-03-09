package main

import (
	"github.com/zephyyrr/plugins"
	"time"
)

type MessangerPlugin struct {
}

func (MessangerPlugin) Name() string {
	return "Messanger"
}

func (MessangerPlugin) Provides() []plugins.Event {
	return []plugins.Event{
		"message",
	}
}

func (MessangerPlugin) Subscribes() []plugins.Event {
	return nil
}

func (MessangerPlugin) Send(e plugins.Event, args plugins.Args) error {
	return nil
}

func (MessangerPlugin) Recieve() (plugins.Event, plugins.Args, error) {
	time.Sleep(time.Second)
	return "message", map[string]interface{}{"message": "Hello from local plugin."}, nil
}
