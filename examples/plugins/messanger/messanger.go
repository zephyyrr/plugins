package main

import (
	"github.com/zephyyrr/plugins"
	"os"
	"time"
)

func main() {
	client, err := plugins.NewClient(plugins.PluginDecl{
		Name: "messanger",
		Provides: []plugins.Event{
			"log",
		},
	}, os.Stdin, os.Stdout)

	if err != nil {
		panic(err)
	}

	go client.Run()

	client.Dispatch("log", plugins.Args{
		"Level":   "INFO",
		"Origin":  "messanger",
		"Message": "Remote pluging through stdin/stdout is connected.",
	})
	time.Sleep(1000 * time.Millisecond)
}
