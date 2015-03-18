package main

import (
	"github.com/zephyyrr/plugins"
	"os"
	"time"
)

func main() {
	client, err := plugins.NewClient(plugins.PluginDecl{
		Name: "pinger",
		Provides: []plugins.Event{
			"ping",
		},
		Subscribes: []plugins.Event{
			"pong",
		},
	}, os.Stdin, os.Stdout)

	if err != nil {
		panic(err)
	}

	client.AddHandler("pong", plugins.HandlerFunc(func(e plugins.Event, args plugins.Args) {
		client.Dispatch("log", plugins.Args{
			"Level":         "INFO",
			"Event":         "ping",
			"Reponse-Event": "pong",
		})
	}))

	go client.Run()

	for {
		client.Dispatch("ping", plugins.Args{
			"message": "Ping from another process",
		})
		time.Sleep(20 * time.Millisecond)
	}
}
