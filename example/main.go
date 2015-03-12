package main

import (
	"fmt"
	"github.com/zephyyrr/plugins"
)

func main() {
	ph := plugins.NewManager()
	ph.Handle(MessangerPlugin{})

	if plugins, err := plugins.LoadAll("./plugins"); err != nil {
		ph.HandleAll(plugins...)
	}

	ph.Muxer().AddHandler("message", plugins.HandlerFunc(func(event plugins.Event, args plugins.Args) {
		fmt.Println(args)
	}))

	ph.ListenAndServe()
}
