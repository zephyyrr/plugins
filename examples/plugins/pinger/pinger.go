package main

import (
	"github.com/zephyyrr/plugins"
	"log"
	"os"
)

const Name = "pinger"

func init() {
	log.SetPrefix(Name)
}

func main() {
	client, err := plugins.NewClient(plugins.PluginDecl{
		Name: Name,
		Provides: []plugins.Event{
			"ping",
			"log",
		},
		Subscribes: []plugins.Event{
			"pong",
		},
	}, os.Stdin, os.Stdout)

	if err != nil {
		panic(err)
	}

	i := 1

	client.AddHandler("pong", plugins.HandlerFunc(func(e plugins.Event, args plugins.Args) {
		err := client.Dispatch("log",
			plugins.Args{
				"Level":      "INFO",
				"Origin":     Name,
				"Message":    "Recieved event",
				"Additional": CallResponse{"ping", e},
			})
		if err != nil {
			log.Panicln(err)
		}
		//time.Sleep(200 * time.Millisecond)

		i++
		client.Dispatch("ping", plugins.Args{
			"message": "Additional ping from another process",
			"count":   i,
		})

	}))

	go client.Run()

	client.Dispatch("ping", plugins.Args{
		"message": "First ping from another process",
	})

	select {}
}

type CallResponse struct {
	In, Out plugins.Event
}
